package admission

import (
	"encoding/json"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
)

// checkReadOnlyRootFilesystem checks if the pod has containers with writable root filesystems
func checkReadOnlyRootFilesystem(request *admissionv1.AdmissionRequest) bool {
	// Parse the Pod object from the request
	pod := &corev1.Pod{}
	err := json.Unmarshal(request.Object.Raw, pod)
	if err != nil {
		return false
	}

	// Check for writable root filesystems
	for _, container := range pod.Spec.Containers {
		if container.SecurityContext != nil && (container.SecurityContext.ReadOnlyRootFilesystem == nil || !*container.SecurityContext.ReadOnlyRootFilesystem) {
			return false // Return false if writable root filesystem is found
		}
	}

	// Passes the check if all containers have read-only root filesystems
	return true
}
