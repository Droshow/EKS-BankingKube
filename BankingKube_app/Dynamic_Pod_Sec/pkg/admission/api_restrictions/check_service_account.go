package api_restrictions

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"

	"gopkg.in/yaml.v2"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
)

var (
	saTracer  = otel.Tracer("bankingkube/dynamicpodsec")
	saMeter   = otel.Meter("bankingkube/dynamicpodsec")
	saDenied  metric.Int64Counter
	saAllowed metric.Int64Counter
)

func init() {
	var err error
	saDenied, err = saMeter.Int64Counter("service_account.denied")
	if err != nil {
		log.Println("Failed to create metric: service_account.denied")
	}
	saAllowed, err = saMeter.Int64Counter("service_account.allowed")
	if err != nil {
		log.Println("Failed to create metric: service_account.allowed")
	}
}

// ServiceAccountRestrictions defines a structure for restricted service accounts
type ServiceAccountRestrictions struct {
	RestrictedServiceAccounts []string `yaml:"restrictedServiceAccounts"`
}

// SecurityPoliciesForSA represents the structure of the security-policies.yaml file
type SecurityPoliciesForSA struct {
	Policies struct {
		ServiceAccountRestrictions ServiceAccountRestrictions `yaml:"serviceAccountRestrictions"`
	} `yaml:"policies"`
}

// CheckServiceAccount ensures that pods do not use restricted service accounts
func CheckServiceAccount(ctx context.Context, request *admissionv1.AdmissionRequest) bool {
	ctx, span := saTracer.Start(ctx, "CheckServiceAccount", trace.WithAttributes(
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
		attribute.String("service_account", pod.Spec.ServiceAccountName),
	)

	// Retrieve the security policies for service account usage
	serviceAccountRestrictions, err := getServiceAccountRestrictions()
	if err != nil {
		log.Println("Failed to load service account restrictions:", err)
		span.SetAttributes(
			attribute.String("error", "failed_to_load_restrictions"),
			attribute.String("result", "denied"),
		)
		span.RecordError(err)
		return false
	}

	// Check if the pod is using a restricted service account
	for _, restrictedAccount := range serviceAccountRestrictions.RestrictedServiceAccounts {
		if pod.Spec.ServiceAccountName == restrictedAccount {
			log.Printf("Pod %s in namespace %s is using a restricted service account: %s\n",
				pod.Name, pod.Namespace, restrictedAccount)

			saDenied.Add(ctx, 1, metric.WithAttributes(
				attribute.String("pod", pod.Name),
				attribute.String("namespace", pod.Namespace),
				attribute.String("service_account", pod.Spec.ServiceAccountName),
				attribute.String("restricted_account", restrictedAccount),
			))

			span.SetAttributes(
				attribute.String("restricted_account", restrictedAccount),
				attribute.String("result", "denied"),
			)

			return false
		}
	}

	// Service account is allowed
	saAllowed.Add(ctx, 1, metric.WithAttributes(
		attribute.String("pod", pod.Name),
		attribute.String("namespace", pod.Namespace),
		attribute.String("service_account", pod.Spec.ServiceAccountName),
	))

	span.SetAttributes(attribute.String("result", "allowed"))
	return true
}

// getServiceAccountRestrictions loads the service account restrictions from the configuration file
func getServiceAccountRestrictions() (*ServiceAccountRestrictions, error) {
	configPath := os.Getenv("SECURITY_POLICIES_PATH")
	if configPath == "" {
		configPath = "configs/security-policies.yaml" // Default path
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var policies SecurityPoliciesForSA
	err = yaml.Unmarshal(data, &policies)
	if err != nil {
		return nil, err
	}

	return &policies.Policies.ServiceAccountRestrictions, nil
}
