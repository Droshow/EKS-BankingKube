package resource_management

import (
	"encoding/json"
	"log"

	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
)

// CheckResourceLimits validates if a pod's containers have resource limits defined
func CheckResourceLimits(request *admissionv1.AdmissionRequest) bool {
	// Parse the Pod object from the request
	pod := &corev1.Pod{}
	err := json.Unmarshal(request.Object.Raw, pod)
	if err != nil {
		log.Println("Failed to parse pod object:", err)
		return false // Fails the validation if the pod can't be parsed
	}

	// Check if the pod's containers have resource limits defined
	for _, container := range pod.Spec.Containers {
		if container.Resources.Limits == nil {
			log.Printf("Pod %s in namespace %s has a container without resource limits: %s\n", pod.Name, pod.Namespace, container.Name)
			return false
		}
	}

	return true // Passes the check if all containers have resource limits defined
}
