package image_security

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
	imgTracer  = otel.Tracer("bankingkube/dynamicpodsec")
	imgMeter   = otel.Meter("bankingkube/dynamicpodsec")
	imgDenied  metric.Int64Counter
	imgAllowed metric.Int64Counter
)

func init() {
	var err error
	imgDenied, err = imgMeter.Int64Counter("image_registry.denied")
	if err != nil {
		log.Println("Failed to create metric: image_registry.denied")
	}
	imgAllowed, err = imgMeter.Int64Counter("image_registry.allowed")
	if err != nil {
		log.Println("Failed to create metric: image_registry.allowed")
	}
}

// ImageSecurity defines a structure for image security policies
type ImageSecurity struct {
	AllowedRegistries []string `yaml:"allowedRegistries"`
}

// SecurityPoliciesImg represents the structure of the security-policies.yaml file
type SecurityPoliciesImg struct {
	Policies struct {
		ImageSecurity ImageSecurity `yaml:"imageSecurity"`
	} `yaml:"policies"`
}

// CheckImageRegistry validates if a pod is using images from allowed registries
func CheckImageRegistry(ctx context.Context, request *admissionv1.AdmissionRequest) bool {
	ctx, span := imgTracer.Start(ctx, "CheckImageRegistry", trace.WithAttributes(
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

	// Retrieve the allowed registries
	imageSecurity, err := getImageSecurity()
	if err != nil {
		log.Println("Failed to load image security policies:", err)
		span.SetAttributes(
			attribute.String("error", "failed_to_load_policies"),
			attribute.String("result", "denied"),
		)
		span.RecordError(err)
		return false
	}

	// Check if the pod's containers are using images from allowed registries
	for _, container := range pod.Spec.Containers {
		span.AddEvent("Checking container image", trace.WithAttributes(
			attribute.String("container", container.Name),
			attribute.String("image", container.Image),
		))

		if !isImageFromAllowedRegistry(container.Image, imageSecurity.AllowedRegistries) {
			log.Printf("Pod %s in namespace %s is using an image from a disallowed registry: %s\n",
				pod.Name, pod.Namespace, container.Image)

			imgDenied.Add(ctx, 1, metric.WithAttributes(
				attribute.String("pod", pod.Name),
				attribute.String("namespace", pod.Namespace),
				attribute.String("container", container.Name),
				attribute.String("image", container.Image),
				attribute.String("denial_reason", "disallowed_registry"),
			))

			span.SetAttributes(
				attribute.String("container", container.Name),
				attribute.String("image", container.Image),
				attribute.String("result", "denied"),
				attribute.String("denial_reason", "disallowed_registry"),
			)

			return false
		}
	}

	// Check init containers as well
	for _, initContainer := range pod.Spec.InitContainers {
		span.AddEvent("Checking init container image", trace.WithAttributes(
			attribute.String("init_container", initContainer.Name),
			attribute.String("image", initContainer.Image),
		))

		if !isImageFromAllowedRegistry(initContainer.Image, imageSecurity.AllowedRegistries) {
			log.Printf("Pod %s in namespace %s is using an init container image from a disallowed registry: %s\n",
				pod.Name, pod.Namespace, initContainer.Image)

			imgDenied.Add(ctx, 1, metric.WithAttributes(
				attribute.String("pod", pod.Name),
				attribute.String("namespace", pod.Namespace),
				attribute.String("init_container", initContainer.Name),
				attribute.String("image", initContainer.Image),
				attribute.String("denial_reason", "disallowed_registry"),
			))

			span.SetAttributes(
				attribute.String("init_container", initContainer.Name),
				attribute.String("image", initContainer.Image),
				attribute.String("result", "denied"),
				attribute.String("denial_reason", "disallowed_registry"),
			)

			return false
		}
	}

	// Check ephemeral containers if they exist
	for _, ephemeralContainer := range pod.Spec.EphemeralContainers {
		span.AddEvent("Checking ephemeral container image", trace.WithAttributes(
			attribute.String("ephemeral_container", ephemeralContainer.Name),
			attribute.String("image", ephemeralContainer.Image),
		))

		if !isImageFromAllowedRegistry(ephemeralContainer.Image, imageSecurity.AllowedRegistries) {
			log.Printf("Pod %s in namespace %s is using an ephemeral container image from a disallowed registry: %s\n",
				pod.Name, pod.Namespace, ephemeralContainer.Image)

			imgDenied.Add(ctx, 1, metric.WithAttributes(
				attribute.String("pod", pod.Name),
				attribute.String("namespace", pod.Namespace),
				attribute.String("ephemeral_container", ephemeralContainer.Name),
				attribute.String("image", ephemeralContainer.Image),
				attribute.String("denial_reason", "disallowed_registry"),
			))

			span.SetAttributes(
				attribute.String("ephemeral_container", ephemeralContainer.Name),
				attribute.String("image", ephemeralContainer.Image),
				attribute.String("result", "denied"),
				attribute.String("denial_reason", "disallowed_registry"),
			)

			return false
		}
	}

	// All images are from allowed registries
	imgAllowed.Add(ctx, 1, metric.WithAttributes(
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

	return true // Passes the check if all images are from allowed registries
}

// getImageSecurity loads the image security policies from the configuration file
func getImageSecurity() (*ImageSecurity, error) {
	configPath := os.Getenv("SECURITY_POLICIES_PATH")
	if configPath == "" {
		configPath = "configs/security-policies.yaml" // Default path
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var policies SecurityPoliciesImg
	err = yaml.Unmarshal(data, &policies)
	if err != nil {
		return nil, err
	}

	return &policies.Policies.ImageSecurity, nil
}

// isImageFromAllowedRegistry checks if an image is from an allowed registry
func isImageFromAllowedRegistry(image string, allowedRegistries []string) bool {
	registry := extractRegistry(image)

	for _, allowedRegistry := range allowedRegistries {
		if strings.HasPrefix(registry, allowedRegistry) {
			return true
		}
	}
	return false
}

// extractRegistry extracts the registry part from an image name
func extractRegistry(image string) string {
	// Handle image names with port number in registry
	// Examples:
	// - docker.io/library/nginx:latest -> docker.io
	// - myregistry.example.com:5000/myapp:1.0 -> myregistry.example.com:5000

	// First split by '/' to get registry part
	parts := strings.SplitN(image, "/", 2)

	if len(parts) == 1 {
		// No slash means it's a Docker Hub official image
		return "docker.io/library"
	}

	// Check if the first part looks like a registry (contains '.' or ':')
	if strings.Contains(parts[0], ".") || strings.Contains(parts[0], ":") {
		return parts[0]
	}

	// Default to Docker Hub
	return "docker.io"
}
