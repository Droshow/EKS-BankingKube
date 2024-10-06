package network_security

import (
	"encoding/json"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	"log"
)

// CheckHostNetwork ensures that pods are not using the host network mode
func CheckHostNetwork(request *admissionv1.AdmissionRequest) bool {
	// Parse the Pod object from the request
	pod := &corev1.Pod{}
	err := json.Unmarshal(request.Object.Raw, pod)
	if err != nil {
		log.Println("Failed to parse Pod object:", err)
		return false // Return false if pod parsing fails
	}

	// Check if hostNetwork is set to true
	if pod.Spec.HostNetwork {
		log.Printf("Pod %s in namespace %s is using host network, which is not allowed.", pod.Name, pod.Namespace)
		return false // Return false if the pod is using host network
	}

	// Passes the check if host network is not used
	return true
}
