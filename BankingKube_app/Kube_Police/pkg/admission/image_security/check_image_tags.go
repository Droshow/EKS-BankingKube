package image_security

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
)

// DisallowedTags defines a structure for disallowed image tags
type DisallowedTags struct {
	DisallowedTags []string `yaml:"disallowedTags"`
}

// CheckImageTags validates if a pod's images are using allowed tags
func CheckImageTags(request *admissionv1.AdmissionRequest) bool {
	// Parse the Pod object from the request
	pod := &corev1.Pod{}
	err := json.Unmarshal(request.Object.Raw, pod)
	if err != nil {
		log.Println("Failed to parse pod object:", err)
		return false // Fails the validation if the pod can't be parsed
	}

	// Retrieve the disallowed tags policies
	disallowedTags, err := getDisallowedTags()
	if err != nil {
		log.Println("Failed to load disallowed tags policies:", err)
		return false
	}

	// Check if the pod's containers are using allowed image tags
	for _, container := range pod.Spec.Containers {
		if !isImageTagAllowed(container.Image, disallowedTags.DisallowedTags) {
			log.Printf("Pod %s in namespace %s is using an image with a disallowed tag: %s\n", pod.Name, pod.Namespace, container.Image)
			return false
		}
	}

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

	var policies struct {
		DisallowedTags DisallowedTags `yaml:"disallowedTags"`
	}

	err = yaml.Unmarshal(data, &policies)
	if err != nil {
		return nil, err
	}

	return &policies.DisallowedTags, nil
}

// isImageTagAllowed checks if an image tag is allowed
func isImageTagAllowed(image string, disallowedTags []string) bool {
	for _, tag := range disallowedTags {
		if strings.HasSuffix(image, ":"+tag) {
			return false
		}
	}
	return true
}
