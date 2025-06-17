package context_capabilities

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
	capTracer  = otel.Tracer("bankingkube/dynamicpodsec")
	capMeter   = otel.Meter("bankingkube/dynamicpodsec")
	capDenied  metric.Int64Counter
	capAllowed metric.Int64Counter
)

func init() {
	var err error
	capDenied, err = capMeter.Int64Counter("capabilities.denied")
	if err != nil {
		log.Println("Failed to create metric: capabilities.denied")
	}
	capAllowed, err = capMeter.Int64Counter("capabilities.allowed")
	if err != nil {
		log.Println("Failed to create metric: capabilities.allowed")
	}
}

// Capabilities defines a structure for capabilities policies
type Capabilities struct {
	AllowedCapabilities    []string `yaml:"allowedCapabilities"`
	DisallowedCapabilities []string `yaml:"disallowedCapabilities"`
	RequiredDrops          []string `yaml:"requiredDrops"`
}

// SecurityPoliciesCap represents the structure of the security-policies.yaml file
type SecurityPoliciesCap struct {
	Policies struct {
		Capabilities Capabilities `yaml:"capabilities"`
	} `yaml:"policies"`
}

// CheckCapabilities ensures that the pod does not use dangerous Linux capabilities
func CheckCapabilities(ctx context.Context, request *admissionv1.AdmissionRequest) bool {
	ctx, span := capTracer.Start(ctx, "CheckCapabilities", trace.WithAttributes(
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

	// Retrieve the capabilities policies
	capabilities, err := getCapabilities()
	if err != nil {
		log.Println("Failed to load capabilities policies:", err)
		span.SetAttributes(
			attribute.String("error", "failed_to_load_policies"),
			attribute.String("result", "denied"),
		)
		span.RecordError(err)
		return false
	}

	// Check capabilities for each container
	for _, container := range pod.Spec.Containers {
		span.AddEvent("Checking container", trace.WithAttributes(
			attribute.String("container", container.Name),
		))

		if container.SecurityContext != nil && container.SecurityContext.Capabilities != nil {
			// Check for disallowed capabilities
			for _, cap := range container.SecurityContext.Capabilities.Add {
				if isDisallowedCapability(string(cap), capabilities.DisallowedCapabilities) {
					log.Println("Disallowed capability found in container:", container.Name, "Capability:", cap)

					capDenied.Add(ctx, 1, metric.WithAttributes(
						attribute.String("pod", pod.Name),
						attribute.String("namespace", pod.Namespace),
						attribute.String("container", container.Name),
						attribute.String("disallowed_capability", string(cap)),
						attribute.String("denial_reason", "disallowed_capability"),
					))

					span.SetAttributes(
						attribute.String("container", container.Name),
						attribute.String("disallowed_capability", string(cap)),
						attribute.String("result", "denied"),
						attribute.String("denial_reason", "disallowed_capability"),
					)

					return false // Return false if any disallowed capability is found
				}
			}

			// Check for required capability drops
			if !hasDroppedAllRequiredCapabilities(container.SecurityContext.Capabilities.Drop, capabilities.RequiredDrops) {
				log.Println("Necessary capabilities not dropped in container:", container.Name)

				missingDrops := getMissingRequiredDrops(container.SecurityContext.Capabilities.Drop, capabilities.RequiredDrops)

				capDenied.Add(ctx, 1, metric.WithAttributes(
					attribute.String("pod", pod.Name),
					attribute.String("namespace", pod.Namespace),
					attribute.String("container", container.Name),
					attribute.String("missing_required_drops", string(missingDrops)),
					attribute.String("denial_reason", "missing_required_drops"),
				))

				span.SetAttributes(
					attribute.String("container", container.Name),
					attribute.String("missing_required_drops", string(missingDrops)),
					attribute.String("result", "denied"),
					attribute.String("denial_reason", "missing_required_drops"),
				)

				return false // Return false if necessary capabilities are not dropped
			}
		}
	}

	// Passes the check if no disallowed capabilities are found and all required capabilities are dropped
	capAllowed.Add(ctx, 1, metric.WithAttributes(
		attribute.String("pod", pod.Name),
		attribute.String("namespace", pod.Namespace),
	))

	span.SetAttributes(attribute.String("result", "allowed"))
	return true
}

// isDisallowedCapability checks if a capability is in the list of disallowed capabilities
func isDisallowedCapability(cap string, disallowedCapabilities []string) bool {
	for _, disallowedCap := range disallowedCapabilities {
		if cap == disallowedCap {
			return true
		}
	}
	return false
}

// hasDroppedAllRequiredCapabilities checks if all required capabilities are dropped
func hasDroppedAllRequiredCapabilities(dropped []corev1.Capability, requiredDrops []string) bool {
	requiredDropsMap := make(map[string]bool)
	for _, cap := range requiredDrops {
		requiredDropsMap[cap] = true
	}

	for _, cap := range dropped {
		delete(requiredDropsMap, string(cap))
	}

	return len(requiredDropsMap) == 0
}

// getMissingRequiredDrops returns a string of capabilities that should have been dropped
func getMissingRequiredDrops(dropped []corev1.Capability, requiredDrops []string) string {
	requiredDropsMap := make(map[string]bool)
	for _, cap := range requiredDrops {
		requiredDropsMap[cap] = true
	}

	for _, cap := range dropped {
		delete(requiredDropsMap, string(cap))
	}

	missing := ""
	for cap := range requiredDropsMap {
		if missing != "" {
			missing += ", "
		}
		missing += cap
	}

	return missing
}

// getCapabilities loads the capabilities policies from the configuration file
func getCapabilities() (*Capabilities, error) {
	configPath := os.Getenv("SECURITY_POLICIES_PATH")
	if configPath == "" {
		configPath = "configs/security-policies.yaml" // Default path
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var policies SecurityPoliciesCap
	err = yaml.Unmarshal(data, &policies)
	if err != nil {
		return nil, err
	}

	return &policies.Policies.Capabilities, nil
}
