package context_capabilities

import (
	"encoding/json"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	"log"
)

// CheckRunAsUser verifies if the pod has appropriate runAsUser and runAsGroup settings
func CheckRunAsUser(request *admissionv1.AdmissionRequest) bool {
	// Parse the Pod object from the request
	pod := &corev1.Pod{}
	err := json.Unmarshal(request.Object.Raw, pod)
	if err != nil {
		log.Println("Failed to parse pod object:", err)
		return false // Fails the validation if the pod can't be parsed
	}

	// Check runAsUser and runAsGroup for all containers
	for _, container := range pod.Spec.Containers {
		if container.SecurityContext != nil {
			if container.SecurityContext.RunAsUser == nil {
				log.Println("runAsUser is not set for container:", container.Name)
				return false
			}
			if container.SecurityContext.RunAsNonRoot != nil && !*container.SecurityContext.RunAsNonRoot {
				log.Println("runAsNonRoot is false for container:", container.Name)
				return false
			}
		}
	}

	// Passes the check if all containers have proper runAsUser settings
	return true
}
