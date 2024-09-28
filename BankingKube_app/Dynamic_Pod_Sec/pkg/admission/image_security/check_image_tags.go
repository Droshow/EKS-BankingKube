package image_security

import (
	"encoding/json"
	"log"
	"strings"

	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
)

// CheckImageTags validates if a pod's images are using allowed tags
func CheckImageTags(request *admissionv1.AdmissionRequest) bool {
	// Parse the Pod object from the request
	pod := &corev1.Pod{}
	err := json.Unmarshal(request.Object.Raw, pod)
	if err != nil {
		log.Println("Failed to parse pod object:", err)
		return false // Fails the validation if the pod can't be parsed
	}

	// Check if the pod's containers are using allowed image tags
	for _, container := range pod.Spec.Containers {
		if !isImageTagAllowed(container.Image) {
			log.Printf("Pod %s in namespace %s is using an image with a disallowed tag: %s\n", pod.Name, pod.Namespace, container.Image)
			return false
		}
	}

	return true // Passes the check if all images have allowed tags
}

// isImageTagAllowed checks if an image tag is allowed
func isImageTagAllowed(image string) bool {
	// Disallow the use of the "latest" tag
	if strings.HasSuffix(image, ":latest") {
		return false
	}
	return true
}
