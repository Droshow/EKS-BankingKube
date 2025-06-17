package image_security

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/exec"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"

	"gopkg.in/yaml.v2"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
)

var (
	signTracer     = otel.Tracer("bankingkube/dynamicpodsec")
	signMeter      = otel.Meter("bankingkube/dynamicpodsec")
	signDenied     metric.Int64Counter
	signAllowed    metric.Int64Counter
	signVerified   metric.Int64Counter
	signUnverified metric.Int64Counter
)

func init() {
	var err error
	signDenied, err = signMeter.Int64Counter("image_signing.denied")
	if err != nil {
		log.Println("Failed to create metric: image_signing.denied")
	}
	signAllowed, err = signMeter.Int64Counter("image_signing.allowed")
	if err != nil {
		log.Println("Failed to create metric: image_signing.allowed")
	}
	signVerified, err = signMeter.Int64Counter("image_signing.verified")
	if err != nil {
		log.Println("Failed to create metric: image_signing.verified")
	}
	signUnverified, err = signMeter.Int64Counter("image_signing.unverified")
	if err != nil {
		log.Println("Failed to create metric: image_signing.unverified")
	}
}

// RequireImageSigning defines a structure for the requireImageSigning policy
type RequireImageSigning struct {
	RequireImageSigning bool `yaml:"requireImageSigning"`
}

// SecurityPoliciesSign represents the structure of the security-policies.yaml file
type SecurityPoliciesSign struct {
	Policies struct {
		ImageSecurity struct {
			RequireImageSigning RequireImageSigning `yaml:"requireImageSigning"`
		} `yaml:"imageSecurity"`
	} `yaml:"policies"`
}

// CheckImageSigning validates if a pod's images are signed and verified
func CheckImageSigning(ctx context.Context, request *admissionv1.AdmissionRequest) bool {
	ctx, span := signTracer.Start(ctx, "CheckImageSigning", trace.WithAttributes(
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

	// Retrieve the requireImageSigning policy
	requireImageSigning, err := getRequireImageSigning()
	if err != nil {
		log.Println("Failed to load requireImageSigning policy:", err)
		span.SetAttributes(
			attribute.String("error", "failed_to_load_policy"),
			attribute.String("result", "denied"),
		)
		span.RecordError(err)
		return false
	}

	// Check if image signing enforcement is enabled
	if !requireImageSigning.RequireImageSigning {
		span.SetAttributes(
			attribute.Bool("image_signing_required", false),
			attribute.String("result", "allowed"),
			attribute.String("reason", "signing_not_required"),
		)

		signAllowed.Add(ctx, 1, metric.WithAttributes(
			attribute.String("pod", pod.Name),
			attribute.String("namespace", pod.Namespace),
			attribute.String("reason", "signing_not_required"),
		))

		return true // Passes the check if enforcement is not enabled
	}

	span.SetAttributes(attribute.Bool("image_signing_required", true))

	// Check regular containers
	for _, container := range pod.Spec.Containers {
		containerCtx, containerSpan := signTracer.Start(ctx, "VerifyContainerImage", trace.WithAttributes(
			attribute.String("container", container.Name),
			attribute.String("image", container.Image),
			attribute.String("container_type", "regular"),
		))

		if !isImageSigned(containerCtx, container.Image) {
			log.Printf("Pod %s in namespace %s is using an unsigned image: %s\n",
				pod.Name, pod.Namespace, container.Image)

			signDenied.Add(ctx, 1, metric.WithAttributes(
				attribute.String("pod", pod.Name),
				attribute.String("namespace", pod.Namespace),
				attribute.String("container", container.Name),
				attribute.String("image", container.Image),
				attribute.String("denial_reason", "unsigned_image"),
				attribute.String("container_type", "regular"),
			))

			containerSpan.SetAttributes(
				attribute.String("result", "denied"),
				attribute.String("reason", "unsigned_image"),
			)
			containerSpan.End()

			span.SetAttributes(
				attribute.String("result", "denied"),
				attribute.String("container", container.Name),
				attribute.String("image", container.Image),
				attribute.String("reason", "unsigned_image"),
				attribute.String("container_type", "regular"),
			)

			return false
		}

		containerSpan.SetAttributes(attribute.String("result", "verified"))
		containerSpan.End()
	}

	// Check init containers
	for _, initContainer := range pod.Spec.InitContainers {
		initCtx, initSpan := signTracer.Start(ctx, "VerifyInitContainerImage", trace.WithAttributes(
			attribute.String("init_container", initContainer.Name),
			attribute.String("image", initContainer.Image),
			attribute.String("container_type", "init"),
		))

		if !isImageSigned(initCtx, initContainer.Image) {
			log.Printf("Pod %s in namespace %s is using an unsigned init container image: %s\n",
				pod.Name, pod.Namespace, initContainer.Image)

			signDenied.Add(ctx, 1, metric.WithAttributes(
				attribute.String("pod", pod.Name),
				attribute.String("namespace", pod.Namespace),
				attribute.String("init_container", initContainer.Name),
				attribute.String("image", initContainer.Image),
				attribute.String("denial_reason", "unsigned_image"),
				attribute.String("container_type", "init"),
			))

			initSpan.SetAttributes(
				attribute.String("result", "denied"),
				attribute.String("reason", "unsigned_image"),
			)
			initSpan.End()

			span.SetAttributes(
				attribute.String("result", "denied"),
				attribute.String("init_container", initContainer.Name),
				attribute.String("image", initContainer.Image),
				attribute.String("reason", "unsigned_image"),
				attribute.String("container_type", "init"),
			)

			return false
		}

		initSpan.SetAttributes(attribute.String("result", "verified"))
		initSpan.End()
	}

	// Check ephemeral containers
	for _, ephemeralContainer := range pod.Spec.EphemeralContainers {
		ephemeralCtx, ephemeralSpan := signTracer.Start(ctx, "VerifyEphemeralContainerImage", trace.WithAttributes(
			attribute.String("ephemeral_container", ephemeralContainer.Name),
			attribute.String("image", ephemeralContainer.Image),
			attribute.String("container_type", "ephemeral"),
		))

		if !isImageSigned(ephemeralCtx, ephemeralContainer.Image) {
			log.Printf("Pod %s in namespace %s is using an unsigned ephemeral container image: %s\n",
				pod.Name, pod.Namespace, ephemeralContainer.Image)

			signDenied.Add(ctx, 1, metric.WithAttributes(
				attribute.String("pod", pod.Name),
				attribute.String("namespace", pod.Namespace),
				attribute.String("ephemeral_container", ephemeralContainer.Name),
				attribute.String("image", ephemeralContainer.Image),
				attribute.String("denial_reason", "unsigned_image"),
				attribute.String("container_type", "ephemeral"),
			))

			ephemeralSpan.SetAttributes(
				attribute.String("result", "denied"),
				attribute.String("reason", "unsigned_image"),
			)
			ephemeralSpan.End()

			span.SetAttributes(
				attribute.String("result", "denied"),
				attribute.String("ephemeral_container", ephemeralContainer.Name),
				attribute.String("image", ephemeralContainer.Image),
				attribute.String("reason", "unsigned_image"),
				attribute.String("container_type", "ephemeral"),
			)

			return false
		}

		ephemeralSpan.SetAttributes(attribute.String("result", "verified"))
		ephemeralSpan.End()
	}

	// All images are signed
	signAllowed.Add(ctx, 1, metric.WithAttributes(
		attribute.String("pod", pod.Name),
		attribute.String("namespace", pod.Namespace),
		attribute.Int("container_count", len(pod.Spec.Containers)),
		attribute.Int("init_container_count", len(pod.Spec.InitContainers)),
		attribute.Int("ephemeral_container_count", len(pod.Spec.EphemeralContainers)),
	))

	span.SetAttributes(
		attribute.String("result", "allowed"),
		attribute.Int("container_count", len(pod.Spec.Containers)),
		attribute.Int("init_container_count", len(pod.Spec.InitContainers)),
		attribute.Int("ephemeral_container_count", len(pod.Spec.EphemeralContainers)),
	)

	return true // Passes the check if all images are signed
}

// getRequireImageSigning loads the requireImageSigning policy from the configuration file
func getRequireImageSigning() (*RequireImageSigning, error) {
	configPath := os.Getenv("SECURITY_POLICIES_PATH")
	if configPath == "" {
		configPath = "configs/security-policies.yaml" // Default path
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var policies SecurityPoliciesSign
	err = yaml.Unmarshal(data, &policies)
	if err != nil {
		return nil, err
	}

	return &policies.Policies.ImageSecurity.RequireImageSigning, nil
}

// isImageSigned checks if the image is signed using cosign
func isImageSigned(ctx context.Context, image string) bool {
	ctx, span := signTracer.Start(ctx, "CosignVerify", trace.WithAttributes(
		attribute.String("image", image),
	))
	defer span.End()

	// Get the public key path from the environment variable
	publicKeyPath := os.Getenv("COSIGN_PUBLIC_KEY_PATH")
	if publicKeyPath == "" {
		log.Println("COSIGN_PUBLIC_KEY_PATH environment variable is not set")
		span.SetAttributes(
			attribute.String("error", "missing_public_key_path"),
			attribute.String("result", "unverified"),
		)

		signUnverified.Add(ctx, 1, metric.WithAttributes(
			attribute.String("image", image),
			attribute.String("reason", "missing_public_key_path"),
		))

		return false
	}

	span.SetAttributes(attribute.String("public_key_path", publicKeyPath))

	// Construct the cosign verify command
	cmd := exec.Command("cosign", "verify", "--key", publicKeyPath, image)

	// Run the command and capture the output
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Failed to verify image signature for %s: %v\nOutput: %s", image, err, string(output))
		span.SetAttributes(
			attribute.String("error", "verification_failed"),
			attribute.String("output", string(output)),
			attribute.String("result", "unverified"),
		)
		span.RecordError(err)

		signUnverified.Add(ctx, 1, metric.WithAttributes(
			attribute.String("image", image),
			attribute.String("reason", "verification_failed"),
		))

		return false
	}

	log.Printf("Successfully verified image signature for %s\n", image)
	span.SetAttributes(
		attribute.String("result", "verified"),
		attribute.String("output", string(output)),
	)

	signVerified.Add(ctx, 1, metric.WithAttributes(
		attribute.String("image", image),
	))

	return true
}
