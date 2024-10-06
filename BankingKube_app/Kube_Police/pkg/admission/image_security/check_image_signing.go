package image_security

import (
	"encoding/json"
	"log"

	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
)

// CheckImageSigning validates if a pod's images are signed and verified
func CheckImageSigning(request *admissionv1.AdmissionRequest) bool {
	// Parse the Pod object from the request
	pod := &corev1.Pod{}
	err := json.Unmarshal(request.Object.Raw, pod)
	if err != nil {
		log.Println("Failed to parse pod object:", err)
		return false // Fails the validation if the pod can't be parsed
	}

	// Check if the pod's containers are using signed images
	for _, container := range pod.Spec.Containers {
		if !isImageSigned(container.Image) {
			log.Printf("Pod %s in namespace %s is using an unsigned image: %s\n", pod.Name, pod.Namespace, container.Image)
			return false
		}
	}

	return true // Passes the check if all images are signed
}

// isImageSigned checks if an image is signed
func isImageSigned(image string) bool {
	// This function would typically check the image signature
	// Here we assume all images are signed for demonstration purposes
	return true
}
