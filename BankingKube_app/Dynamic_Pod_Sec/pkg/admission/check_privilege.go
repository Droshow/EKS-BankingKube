package admission

import (
	"encoding/json"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
)

// checkPrivilegedContainers checks if the pod has any privileged containers
func checkPrivilegedContainers(request *admissionv1.AdmissionRequest) bool {
	// Parse the Pod object from the request
	pod := &corev1.Pod{}
	err := json.Unmarshal(request.Object.Raw, pod)
	if err != nil {
		return false // Fails the validation if the pod can't be parsed
	}

	// Check for privileged containers
	for _, container := range pod.Spec.Containers {
		if container.SecurityContext != nil && *container.SecurityContext.Privileged {
			return false // Return false if any container is privileged
		}
	}

	// Passes the check if no privileged containers are found
	return true
}
