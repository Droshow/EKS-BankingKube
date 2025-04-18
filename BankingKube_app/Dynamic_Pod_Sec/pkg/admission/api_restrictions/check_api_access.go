package api_restrictions

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"

	"gopkg.in/yaml.v2"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
)

var (
	tracer     = otel.Tracer("bankingkube/dynamicpodsec")
	meter      = otel.Meter("bankingkube/dynamicpodsec")
	apiDenied  metric.Int64Counter
	apiAllowed metric.Int64Counter
)

func init() {
	var err error
	apiDenied, err = meter.Int64Counter("pod_api_access.denied")
	if err != nil {
		log.Println("Failed to create metric: pod_api_access.denied")
	}
	apiAllowed, err = meter.Int64Counter("pod_api_access.allowed")
	if err != nil {
		log.Println("Failed to create metric: pod_api_access.allowed")
	}
}

// APIRestrictions defines a structure for restricted API paths
type APIRestrictions struct {
	RestrictedAPIPaths []string `yaml:"restrictedAPIPaths"`
}

// SecurityPolicies represents the structure of the security-policies.yaml file
type SecurityPolicies struct {
	Policies struct {
		APIRestrictions APIRestrictions `yaml:"apiRestrictions"`
	} `yaml:"policies"`
}

// CheckAPIAccess validates if a pod is trying to access restricted Kubernetes API paths
func CheckAPIAccess(ctx context.Context, request *admissionv1.AdmissionRequest) bool {
	ctx, span := tracer.Start(ctx, "CheckAPIAccess", trace.WithAttributes(
		attribute.String("operation", string(request.Operation)),
		attribute.String("resource", request.Resource.Resource),
	))
	defer span.End()

	// Parse the Pod object from the request
	pod := &corev1.Pod{}
	err := json.Unmarshal(request.Object.Raw, pod)
	if err != nil {
		log.Println("Failed to parse pod object:", err)
		span.SetAttributes(
			attribute.String("error", "failed_to_parse_pod"),
			attribute.String("result", "denied"),
		)
		span.RecordError(err)
		return false // Fails the validation if the pod can't be parsed
	}

	span.SetAttributes(
		attribute.String("pod", pod.Name),
		attribute.String("namespace", pod.Namespace),
	)

	// Retrieve the API restrictions
	apiRestrictions, err := getAPIRestrictions()
	if err != nil {
		log.Println("Failed to load API restrictions:", err)
		span.SetAttributes(
			attribute.String("error", "failed_to_load_restrictions"),
			attribute.String("result", "denied"),
		)
		span.RecordError(err)
		return false
	}

	// Check if the pod is trying to access restricted API paths
	requestPath := "/" + request.Resource.Group + "/" + request.Resource.Version + "/" + request.Resource.Resource + "/" + request.Name
	for _, restrictedPath := range apiRestrictions.RestrictedAPIPaths {
		// Convert wildcards to proper regex patterns for matching
		pattern := strings.ReplaceAll(restrictedPath, "*", ".*")
		if strings.Contains(requestPath, pattern) {
			log.Printf("Pod %s in namespace %s is attempting to access restricted API path: %s\n",
				pod.Name, pod.Namespace, restrictedPath)

			apiDenied.Add(ctx, 1, metric.WithAttributes(
				attribute.String("pod", pod.Name),
				attribute.String("namespace", pod.Namespace),
				attribute.String("restricted_path", restrictedPath),
				attribute.String("request_path", requestPath),
			))

			span.SetAttributes(
				attribute.String("restricted_path", restrictedPath),
				attribute.String("request_path", requestPath),
				attribute.String("result", "denied"),
			)

			return false
		}
	}

	// Request is allowed
	apiAllowed.Add(ctx, 1, metric.WithAttributes(
		attribute.String("pod", pod.Name),
		attribute.String("namespace", pod.Namespace),
		attribute.String("request_path", requestPath),
	))

	span.SetAttributes(attribute.String("result", "allowed"))
	return true
}

// getAPIRestrictions loads the API restrictions from the configuration file
func getAPIRestrictions() (*APIRestrictions, error) {
	configPath := os.Getenv("SECURITY_POLICIES_PATH")
	if configPath == "" {
		configPath = "configs/security-policies.yaml" // Default path
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var policies SecurityPolicies
	err = yaml.Unmarshal(data, &policies)
	if err != nil {
		return nil, err
	}

	return &policies.Policies.APIRestrictions, nil
}
