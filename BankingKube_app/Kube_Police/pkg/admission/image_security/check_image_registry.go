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

// ImageSecurity defines a structure for image security policies
type ImageSecurity struct {
	AllowedRegistries []string `yaml:"allowedRegistries"`
}

// CheckImageRegistry validates if a pod is using images from allowed registries
func CheckImageRegistry(request *admissionv1.AdmissionRequest) bool {
	// Parse the Pod object from the request
	pod := &corev1.Pod{}
	err := json.Unmarshal(request.Object.Raw, pod)
	if err != nil {
		log.Println("Failed to parse pod object:", err)
		return false // Fails the validation if the pod can't be parsed
	}

	// Retrieve the allowed registries
	imageSecurity, err := getImageSecurity()
	if err != nil {
		log.Println("Failed to load image security policies:", err)
		return false
	}

	// Check if the pod's containers are using images from allowed registries
	for _, container := range pod.Spec.Containers {
		if !isImageFromAllowedRegistry(container.Image, imageSecurity.AllowedRegistries) {
			log.Printf("Pod %s in namespace %s is using an image from a disallowed registry: %s\n", pod.Name, pod.Namespace, container.Image)
			return false
		}
	}

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

	var policies struct {
		ImageSecurity ImageSecurity `yaml:"imageSecurity"`
	}

	err = yaml.Unmarshal(data, &policies)
	if err != nil {
		return nil, err
	}

	return &policies.ImageSecurity, nil
}

// isImageFromAllowedRegistry checks if an image is from an allowed registry
func isImageFromAllowedRegistry(image string, allowedRegistries []string) bool {
	for _, registry := range allowedRegistries {
		if strings.HasPrefix(image, registry) {
			return true
		}
	}
	return false
}
