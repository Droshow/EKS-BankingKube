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
	tagTracer  = otel.Tracer("bankingkube/dynamicpodsec")
	tagMeter   = otel.Meter("bankingkube/dynamicpodsec")
	tagDenied  metric.Int64Counter
	tagAllowed metric.Int64Counter
)

func init() {
	var err error
	tagDenied, err = tagMeter.Int64Counter("image_tags.denied")
	if err != nil {
		log.Println("Failed to create metric: image_tags.denied")
	}
	tagAllowed, err = tagMeter.Int64Counter("image_tags.allowed")
	if err != nil {
		log.Println("Failed to create metric: image_tags.allowed")
	}
}

// DisallowedTags defines a structure for disallowed image tags
type DisallowedTags struct {
	DisallowedTags []string `yaml:"disallowedTags"`
}

// SecurityPoliciesTags represents the structure of the security-policies.yaml file
type SecurityPoliciesTags struct {
	Policies struct {
		ImageSecurity struct {
			DisallowedTags DisallowedTags `yaml:"disallowedTags"`
		} `yaml:"imageSecurity"`
	} `yaml:"policies"`
}

// CheckImageTags validates if a pod's images are using allowed tags
func CheckImageTags(ctx context.Context, request *admissionv1.AdmissionRequest) bool {
	ctx, span := tagTracer.Start(ctx, "CheckImageTags", trace.WithAttributes(
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

	// Retrieve the disallowed tags policies
	disallowedTags, err := getDisallowedTags()
	if err != nil {
		log.Println("Failed to load disallowed tags policies:", err)
		span.SetAttributes(
			attribute.String("error", "failed_to_load_policies"),
			attribute.String("result", "denied"),
		)
		span.RecordError(err)
		return false
	}

	// Check regular containers
	for _, container := range pod.Spec.Containers {
		span.AddEvent("Checking container image tag", trace.WithAttributes(
			attribute.String("container", container.Name),
			attribute.String("image", container.Image),
		))

		tag := extractImageTag(container.Image)
		if !isImageTagAllowed(container.Image, disallowedTags.DisallowedTags) {
			log.Printf("Pod %s in namespace %s is using an image with a disallowed tag: %s\n",
				pod.Name, pod.Namespace, container.Image)

			tagDenied.Add(ctx, 1, metric.WithAttributes(
				attribute.String("pod", pod.Name),
				attribute.String("namespace", pod.Namespace),
				attribute.String("container", container.Name),
				attribute.String("image", container.Image),
				attribute.String("tag", tag),
				attribute.String("container_type", "regular"),
				attribute.String("denial_reason", "disallowed_tag"),
			))

			span.SetAttributes(
				attribute.String("container", container.Name),
				attribute.String("image", container.Image),
				attribute.String("tag", tag),
				attribute.String("result", "denied"),
				attribute.String("denial_reason", "disallowed_tag"),
			)

			return false
		}
	}

	// Check init containers
	for _, initContainer := range pod.Spec.InitContainers {
		span.AddEvent("Checking init container image tag", trace.WithAttributes(
			attribute.String("init_container", initContainer.Name),
			attribute.String("image", initContainer.Image),
		))

		tag := extractImageTag(initContainer.Image)
		if !isImageTagAllowed(initContainer.Image, disallowedTags.DisallowedTags) {
			log.Printf("Pod %s in namespace %s is using an init container image with a disallowed tag: %s\n",
				pod.Name, pod.Namespace, initContainer.Image)

			tagDenied.Add(ctx, 1, metric.WithAttributes(
				attribute.String("pod", pod.Name),
				attribute.String("namespace", pod.Namespace),
				attribute.String("init_container", initContainer.Name),
				attribute.String("image", initContainer.Image),
				attribute.String("tag", tag),
				attribute.String("container_type", "init"),
				attribute.String("denial_reason", "disallowed_tag"),
			))

			span.SetAttributes(
				attribute.String("init_container", initContainer.Name),
				attribute.String("image", initContainer.Image),
				attribute.String("tag", tag),
				attribute.String("result", "denied"),
				attribute.String("denial_reason", "disallowed_tag"),
			)

			return false
		}
	}

	// Check ephemeral containers
	for _, ephemeralContainer := range pod.Spec.EphemeralContainers {
		span.AddEvent("Checking ephemeral container image tag", trace.WithAttributes(
			attribute.String("ephemeral_container", ephemeralContainer.Name),
			attribute.String("image", ephemeralContainer.Image),
		))

		tag := extractImageTag(ephemeralContainer.Image)
		if !isImageTagAllowed(ephemeralContainer.Image, disallowedTags.DisallowedTags) {
			log.Printf("Pod %s in namespace %s is using an ephemeral container image with a disallowed tag: %s\n",
				pod.Name, pod.Namespace, ephemeralContainer.Image)

			tagDenied.Add(ctx, 1, metric.WithAttributes(
				attribute.String("pod", pod.Name),
				attribute.String("namespace", pod.Namespace),
				attribute.String("ephemeral_container", ephemeralContainer.Name),
				attribute.String("image", ephemeralContainer.Image),
				attribute.String("tag", tag),
				attribute.String("container_type", "ephemeral"),
				attribute.String("denial_reason", "disallowed_tag"),
			))

			span.SetAttributes(
				attribute.String("ephemeral_container", ephemeralContainer.Name),
				attribute.String("image", ephemeralContainer.Image),
				attribute.String("tag", tag),
				attribute.String("result", "denied"),
				attribute.String("denial_reason", "disallowed_tag"),
			)

			return false
		}
	}

	// All images have allowed tags
	tagAllowed.Add(ctx, 1, metric.WithAttributes(
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

	return true // Passes the check if all images have allowed tags
}

// getDisallowedTags loads the disallowed tags policies from the configuration file
func getDisallowedTags() (*DisallowedTags, error) {
	configPath := os.Getenv("SECURITY_POLICIES_PATH")
	if configPath == "" {
		configPath = "configs/security-policies.yaml" // Default path
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var policies SecurityPoliciesTags
	err = yaml.Unmarshal(data, &policies)
	if err != nil {
		return nil, err
	}

	return &policies.Policies.ImageSecurity.DisallowedTags, nil
}

// isImageTagAllowed checks if an image tag is allowed
func isImageTagAllowed(image string, disallowedTags []string) bool {
	tag := extractImageTag(image)

	for _, disallowedTag := range disallowedTags {
		if tag == disallowedTag {
			return false
		}
	}
	return true
}

// extractImageTag extracts the tag part from an image name
func extractImageTag(image string) string {
	// Default tag if none is specified
	defaultTag := "latest"

	// First handle digest format (image@sha256:1234...)
	if strings.Contains(image, "@") {
		return "digest" // Special case for digest references
	}

	// Handle normal tag format (image:tag)
	parts := strings.Split(image, ":")
	if len(parts) == 1 {
		return defaultTag
	}

	// Get the last part which should be the tag
	// This handles cases like host:port/image:tag correctly
	return parts[len(parts)-1]
}
