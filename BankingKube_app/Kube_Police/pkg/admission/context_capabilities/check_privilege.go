package context_capabilities

import (
	"encoding/json"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	"log"
)

// checkPrivilegedContainers checks if the pod has any privileged containers or init containers
func CheckPrivilegedContainers(request *admissionv1.AdmissionRequest) bool {
	// Parse the Pod object from the request
	pod := &corev1.Pod{}
	err := json.Unmarshal(request.Object.Raw, pod)
	if err != nil {
		log.Println("Failed to parse pod object:", err)
		return false // Fails the validation if the pod can't be parsed
	}

	// Check for privileged containers
	for _, container := range pod.Spec.Containers {
		if container.SecurityContext == nil {
			log.Println("Warning: Container", container.Name, "does not have a SecurityContext defined")
			return false // Consider it a security violation if SecurityContext is not defined
		}
		if container.SecurityContext.Privileged != nil && *container.SecurityContext.Privileged {
			return false // Return false if any container is privileged
		}
	}

	// Check for privileged init containers
	for _, initContainer := range pod.Spec.InitContainers {
		if initContainer.SecurityContext == nil {
			log.Println("Warning: Init container", initContainer.Name, "does not have a SecurityContext defined")
			return false // Consider it a security violation if SecurityContext is not defined
		}
		if initContainer.SecurityContext.Privileged != nil && *initContainer.SecurityContext.Privileged {
			return false // Return false if any init container is privileged
		}
	}

	// Passes the check if no privileged containers or init containers are found
	return true
}
