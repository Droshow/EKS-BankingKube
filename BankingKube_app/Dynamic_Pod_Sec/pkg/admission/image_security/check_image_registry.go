package image_security

import (
	"encoding/json"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	"log"
	"strings"
)

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
	allowedRegistries := getAllowedRegistries()

	// Check if the pod's containers are using images from allowed registries
	for _, container := range pod.Spec.Containers {
		if !isImageFromAllowedRegistry(container.Image, allowedRegistries) {
			log.Printf("Pod %s in namespace %s is using an image from a disallowed registry: %s\n", pod.Name, pod.Namespace, container.Image)
			return false
		}
	}

	return true // Passes the check if all images are from allowed registries
}

// getAllowedRegistries loads the allowed registries from the configuration file
func getAllowedRegistries() []string {
	// This function would typically load allowed registries from a config file like security-policies.yaml
	// Here we define them statically for demonstration purposes
	return []string{
		"myregistry.com",
		"trustedregistry.com",
	}
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
