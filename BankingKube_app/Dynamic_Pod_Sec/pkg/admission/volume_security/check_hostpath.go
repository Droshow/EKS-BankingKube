package volume_security

import (
	"encoding/json"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
)

// checkHostPath checks if the pod has any hostPath volumes
func CheckHostPath(request *admissionv1.AdmissionRequest) bool {
	// Parse the Pod object from the request
	pod := &corev1.Pod{}
	err := json.Unmarshal(request.Object.Raw, pod)
	if err != nil {
		return false
	}

	// Check for hostPath volumes
	for _, volume := range pod.Spec.Volumes {
		if volume.HostPath != nil {
			return false // Return false if any hostPath volume is found
		}
	}

	// Passes the check if no hostPath volumes are found
	return true
}
