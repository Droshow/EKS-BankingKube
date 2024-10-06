package context_capabilities

import (
	"encoding/json"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	"log"
)

// CheckCapabilities ensures that the pod does not use dangerous Linux capabilities
func CheckCapabilities(request *admissionv1.AdmissionRequest) bool {
	// Parse the Pod object from the request
	pod := &corev1.Pod{}
	err := json.Unmarshal(request.Object.Raw, pod)
	if err != nil {
		log.Println("Failed to parse pod object:", err)
		return false // Fails the validation if the pod can't be parsed
	}

	// Check capabilities for each container
	for _, container := range pod.Spec.Containers {
		if container.SecurityContext != nil && container.SecurityContext.Capabilities != nil {
			for _, cap := range container.SecurityContext.Capabilities.Add {
				if IsDangerousCapability(string(cap)) {
					log.Println("Dangerous capability found in container:", container.Name, "Capability:", cap)
					return false // Return false if any dangerous capability is found
				}
			}
		}
	}

	// Passes the check if no dangerous capabilities are found
	return true
}

// IsDangerousCapability checks if a capability is considered dangerous
func IsDangerousCapability(cap string) bool {
	dangerousCapabilities := map[string]bool{
		"CAP_SYS_ADMIN":  true,
		"CAP_NET_ADMIN":  true,
		"CAP_SYS_MODULE": true,
		"CAP_SYS_PTRACE": true,
	}

	return dangerousCapabilities[cap]
}
