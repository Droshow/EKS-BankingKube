package context_capabilities

import (
	"encoding/json"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	"log"
)

// ValidateCapability ensures that dangerous capabilities are not added and necessary ones are dropped
func ValidateCapability(request *admissionv1.AdmissionRequest) bool {
	// Parse the Pod object from the request
	pod := &corev1.Pod{}
	err := json.Unmarshal(request.Object.Raw, pod)
	if err != nil {
		log.Println("Failed to parse pod object:", err)
		return false // Fails the validation if the pod can't be parsed
	}

	// Check and ensure necessary capabilities are dropped
	for _, container := range pod.Spec.Containers {
		if container.SecurityContext != nil && container.SecurityContext.Capabilities != nil {
			if !HasDroppedAllCapabilities(container.SecurityContext.Capabilities.Drop) {
				log.Println("Necessary capabilities not dropped in container:", container.Name)
				return false // Return false if necessary capabilities are not dropped
			}
		}
	}

	// Passes the check if all required capabilities are dropped
	return true
}

// HasDroppedAllCapabilities checks if all required capabilities are dropped
func HasDroppedAllCapabilities(dropped []corev1.Capability) bool {
	requiredDrops := map[string]bool{
		"CAP_SYS_ADMIN": true,
		"CAP_NET_ADMIN": true,
	}

	for _, cap := range dropped {
		delete(requiredDrops, string(cap))
	}

	return len(requiredDrops) == 0
}
