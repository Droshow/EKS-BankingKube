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
	pscTracer  = otel.Tracer("bankingkube/dynamicpodsec")
	pscMeter   = otel.Meter("bankingkube/dynamicpodsec")
	pscDenied  metric.Int64Counter
	pscAllowed metric.Int64Counter
)

func init() {
	var err error
	pscDenied, err = pscMeter.Int64Counter("pod_security_context.denied")
	if err != nil {
		log.Println("Failed to create metric: pod_security_context.denied")
	}
	pscAllowed, err = pscMeter.Int64Counter("pod_security_context.allowed")
	if err != nil {
		log.Println("Failed to create metric: pod_security_context.allowed")
	}
}

// PodSecurityContext defines a structure for pod security context policies
type PodSecurityContext struct {
	AllowPrivilegeEscalation bool `yaml:"allowPrivilegeEscalation"`
	RunAsNonRoot             bool `yaml:"runAsNonRoot"`
	ReadOnlyRootFilesystem   bool `yaml:"readOnlyRootFilesystem"`
}

// SecurityPoliciesPSC represents the structure of the security-policies.yaml file
type SecurityPoliciesPSC struct {
	Policies struct {
		PodSecurityContext PodSecurityContext `yaml:"podSecurityContext"`
	} `yaml:"policies"`
}

// CheckPodSecurityContext checks if the pod complies with the pod security context policies
func CheckPodSecurityContext(ctx context.Context, request *admissionv1.AdmissionRequest) bool {
	ctx, span := pscTracer.Start(ctx, "CheckPodSecurityContext", trace.WithAttributes(
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

	// Retrieve the security policies for pod security context
	podSecurityContext, err := getPodSecurityContext()
	if err != nil {
		log.Println("Failed to load pod security context policies:", err)
		span.SetAttributes(
			attribute.String("error", "failed_to_load_policies"),
			attribute.String("result", "denied"),
		)
		span.RecordError(err)
		return false
	}

	// Check for privileged containers
	for _, container := range pod.Spec.Containers {
		span.AddEvent("Checking container security context", trace.WithAttributes(
			attribute.String("container", container.Name),
			attribute.String("container_type", "regular"),
		))

		if container.SecurityContext == nil {
			log.Println("Warning: Container", container.Name, "does not have a SecurityContext defined")

			pscDenied.Add(ctx, 1, metric.WithAttributes(
				attribute.String("pod", pod.Name),
				attribute.String("namespace", pod.Namespace),
				attribute.String("container", container.Name),
				attribute.String("denial_reason", "missing_security_context"),
			))

			span.SetAttributes(
				attribute.String("container", container.Name),
				attribute.String("result", "denied"),
				attribute.String("denial_reason", "missing_security_context"),
			)

			return false // Consider it a security violation if SecurityContext is not defined
		}

		if container.SecurityContext.Privileged != nil && *container.SecurityContext.Privileged {
			log.Println("Container", container.Name, "is privileged, which is not allowed")

			pscDenied.Add(ctx, 1, metric.WithAttributes(
				attribute.String("pod", pod.Name),
				attribute.String("namespace", pod.Namespace),
				attribute.String("container", container.Name),
				attribute.String("denial_reason", "privileged_container"),
			))

			span.SetAttributes(
				attribute.String("container", container.Name),
				attribute.String("result", "denied"),
				attribute.String("denial_reason", "privileged_container"),
			)

			return false // Return false if any container is privileged
		}

		if container.SecurityContext.AllowPrivilegeEscalation != nil &&
			*container.SecurityContext.AllowPrivilegeEscalation != podSecurityContext.AllowPrivilegeEscalation {
			log.Println("Container", container.Name, "has AllowPrivilegeEscalation set to",
				*container.SecurityContext.AllowPrivilegeEscalation, "which does not match the policy")

			pscDenied.Add(ctx, 1, metric.WithAttributes(
				attribute.String("pod", pod.Name),
				attribute.String("namespace", pod.Namespace),
				attribute.String("container", container.Name),
				attribute.String("denial_reason", "invalid_privilege_escalation"),
				attribute.Bool("set_value", *container.SecurityContext.AllowPrivilegeEscalation),
				attribute.Bool("required_value", podSecurityContext.AllowPrivilegeEscalation),
			))

			span.SetAttributes(
				attribute.String("container", container.Name),
				attribute.String("result", "denied"),
				attribute.String("denial_reason", "invalid_privilege_escalation"),
				attribute.Bool("set_value", *container.SecurityContext.AllowPrivilegeEscalation),
				attribute.Bool("required_value", podSecurityContext.AllowPrivilegeEscalation),
			)

			return false
		}

		if container.SecurityContext.RunAsNonRoot != nil &&
			*container.SecurityContext.RunAsNonRoot != podSecurityContext.RunAsNonRoot {
			log.Println("Container", container.Name, "has RunAsNonRoot set to",
				*container.SecurityContext.RunAsNonRoot, "which does not match the policy")

			pscDenied.Add(ctx, 1, metric.WithAttributes(
				attribute.String("pod", pod.Name),
				attribute.String("namespace", pod.Namespace),
				attribute.String("container", container.Name),
				attribute.String("denial_reason", "invalid_run_as_non_root"),
				attribute.Bool("set_value", *container.SecurityContext.RunAsNonRoot),
				attribute.Bool("required_value", podSecurityContext.RunAsNonRoot),
			))

			span.SetAttributes(
				attribute.String("container", container.Name),
				attribute.String("result", "denied"),
				attribute.String("denial_reason", "invalid_run_as_non_root"),
				attribute.Bool("set_value", *container.SecurityContext.RunAsNonRoot),
				attribute.Bool("required_value", podSecurityContext.RunAsNonRoot),
			)

			return false
		}

		if container.SecurityContext.ReadOnlyRootFilesystem != nil &&
			*container.SecurityContext.ReadOnlyRootFilesystem != podSecurityContext.ReadOnlyRootFilesystem {
			log.Println("Container", container.Name, "has ReadOnlyRootFilesystem set to",
				*container.SecurityContext.ReadOnlyRootFilesystem, "which does not match the policy")

			pscDenied.Add(ctx, 1, metric.WithAttributes(
				attribute.String("pod", pod.Name),
				attribute.String("namespace", pod.Namespace),
				attribute.String("container", container.Name),
				attribute.String("denial_reason", "invalid_read_only_root_fs"),
				attribute.Bool("set_value", *container.SecurityContext.ReadOnlyRootFilesystem),
				attribute.Bool("required_value", podSecurityContext.ReadOnlyRootFilesystem),
			))

			span.SetAttributes(
				attribute.String("container", container.Name),
				attribute.String("result", "denied"),
				attribute.String("denial_reason", "invalid_read_only_root_fs"),
				attribute.Bool("set_value", *container.SecurityContext.ReadOnlyRootFilesystem),
				attribute.Bool("required_value", podSecurityContext.ReadOnlyRootFilesystem),
			)

			return false
		}
	}

	// Check for privileged init containers
	for _, initContainer := range pod.Spec.InitContainers {
		span.AddEvent("Checking init container security context", trace.WithAttributes(
			attribute.String("init_container", initContainer.Name),
			attribute.String("container_type", "init"),
		))

		if initContainer.SecurityContext == nil {
			log.Println("Warning: Init container", initContainer.Name, "does not have a SecurityContext defined")

			pscDenied.Add(ctx, 1, metric.WithAttributes(
				attribute.String("pod", pod.Name),
				attribute.String("namespace", pod.Namespace),
				attribute.String("init_container", initContainer.Name),
				attribute.String("denial_reason", "missing_security_context"),
			))

			span.SetAttributes(
				attribute.String("init_container", initContainer.Name),
				attribute.String("result", "denied"),
				attribute.String("denial_reason", "missing_security_context"),
			)

			return false
		}

		if initContainer.SecurityContext.Privileged != nil && *initContainer.SecurityContext.Privileged {
			log.Println("Init container", initContainer.Name, "is privileged, which is not allowed")

			pscDenied.Add(ctx, 1, metric.WithAttributes(
				attribute.String("pod", pod.Name),
				attribute.String("namespace", pod.Namespace),
				attribute.String("init_container", initContainer.Name),
				attribute.String("denial_reason", "privileged_container"),
			))

			span.SetAttributes(
				attribute.String("init_container", initContainer.Name),
				attribute.String("result", "denied"),
				attribute.String("denial_reason", "privileged_container"),
			)

			return false
		}

		if initContainer.SecurityContext.AllowPrivilegeEscalation != nil &&
			*initContainer.SecurityContext.AllowPrivilegeEscalation != podSecurityContext.AllowPrivilegeEscalation {
			log.Println("Init container", initContainer.Name, "has AllowPrivilegeEscalation set to",
				*initContainer.SecurityContext.AllowPrivilegeEscalation, "which does not match the policy")

			pscDenied.Add(ctx, 1, metric.WithAttributes(
				attribute.String("pod", pod.Name),
				attribute.String("namespace", pod.Namespace),
				attribute.String("init_container", initContainer.Name),
				attribute.String("denial_reason", "invalid_privilege_escalation"),
				attribute.Bool("set_value", *initContainer.SecurityContext.AllowPrivilegeEscalation),
				attribute.Bool("required_value", podSecurityContext.AllowPrivilegeEscalation),
			))

			span.SetAttributes(
				attribute.String("init_container", initContainer.Name),
				attribute.String("result", "denied"),
				attribute.String("denial_reason", "invalid_privilege_escalation"),
				attribute.Bool("set_value", *initContainer.SecurityContext.AllowPrivilegeEscalation),
				attribute.Bool("required_value", podSecurityContext.AllowPrivilegeEscalation),
			)

			return false
		}

		if initContainer.SecurityContext.RunAsNonRoot != nil &&
			*initContainer.SecurityContext.RunAsNonRoot != podSecurityContext.RunAsNonRoot {
			log.Println("Init container", initContainer.Name, "has RunAsNonRoot set to",
				*initContainer.SecurityContext.RunAsNonRoot, "which does not match the policy")

			pscDenied.Add(ctx, 1, metric.WithAttributes(
				attribute.String("pod", pod.Name),
				attribute.String("namespace", pod.Namespace),
				attribute.String("init_container", initContainer.Name),
				attribute.String("denial_reason", "invalid_run_as_non_root"),
				attribute.Bool("set_value", *initContainer.SecurityContext.RunAsNonRoot),
				attribute.Bool("required_value", podSecurityContext.RunAsNonRoot),
			))

			span.SetAttributes(
				attribute.String("init_container", initContainer.Name),
				attribute.String("result", "denied"),
				attribute.String("denial_reason", "invalid_run_as_non_root"),
				attribute.Bool("set_value", *initContainer.SecurityContext.RunAsNonRoot),
				attribute.Bool("required_value", podSecurityContext.RunAsNonRoot),
			)

			return false
		}

		if initContainer.SecurityContext.ReadOnlyRootFilesystem != nil &&
			*initContainer.SecurityContext.ReadOnlyRootFilesystem != podSecurityContext.ReadOnlyRootFilesystem {
			log.Println("Init container", initContainer.Name, "has ReadOnlyRootFilesystem set to",
				*initContainer.SecurityContext.ReadOnlyRootFilesystem, "which does not match the policy")

			pscDenied.Add(ctx, 1, metric.WithAttributes(
				attribute.String("pod", pod.Name),
				attribute.String("namespace", pod.Namespace),
				attribute.String("init_container", initContainer.Name),
				attribute.String("denial_reason", "invalid_read_only_root_fs"),
				attribute.Bool("set_value", *initContainer.SecurityContext.ReadOnlyRootFilesystem),
				attribute.Bool("required_value", podSecurityContext.ReadOnlyRootFilesystem),
			))

			span.SetAttributes(
				attribute.String("init_container", initContainer.Name),
				attribute.String("result", "denied"),
				attribute.String("denial_reason", "invalid_read_only_root_fs"),
				attribute.Bool("set_value", *initContainer.SecurityContext.ReadOnlyRootFilesystem),
				attribute.Bool("required_value", podSecurityContext.ReadOnlyRootFilesystem),
			)

			return false
		}
	}

	// Passes the check if all containers comply with security policies
	pscAllowed.Add(ctx, 1, metric.WithAttributes(
		attribute.String("pod", pod.Name),
		attribute.String("namespace", pod.Namespace),
		attribute.Int("container_count", len(pod.Spec.Containers)),
		attribute.Int("init_container_count", len(pod.Spec.InitContainers)),
	))

	span.SetAttributes(
		attribute.String("result", "allowed"),
		attribute.Int("container_count", len(pod.Spec.Containers)),
		attribute.Int("init_container_count", len(pod.Spec.InitContainers)),
	)

	return true
}

// getPodSecurityContext loads the pod security context policies from the configuration file
func getPodSecurityContext() (*PodSecurityContext, error) {
	configPath := os.Getenv("SECURITY_POLICIES_PATH")
	if configPath == "" {
		configPath = "configs/security-policies.yaml" // Default path
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var policies SecurityPoliciesPSC
	err = yaml.Unmarshal(data, &policies)
	if err != nil {
		return nil, err
	}

	return &policies.Policies.PodSecurityContext, nil
}
