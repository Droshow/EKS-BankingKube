package api_restrictions

import (
	"encoding/json"
	"log"
	"strings"

	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
)

// CheckAPIAccess validates if a pod is trying to access restricted Kubernetes API paths
func CheckAPIAccess(request *admissionv1.AdmissionRequest) bool {
	// Parse the Pod object from the request
	pod := &corev1.Pod{}
	err := json.Unmarshal(request.Object.Raw, pod)
	if err != nil {
		log.Println("Failed to parse pod object:", err)
		return false // Fails the validation if the pod can't be parsed
	}

	// Retrieve the security policies for API access
	apiRestrictions := getAPIRestrictions()

	// Check if the pod's service account is allowed to access restricted API paths
	for _, restrictedPath := range apiRestrictions.RestrictedAPIPaths {
		if strings.Contains(request.Name, restrictedPath) || strings.Contains(request.Namespace, restrictedPath) {
			log.Printf("Pod %s in namespace %s is attempting to access restricted API path: %s\n", pod.Name, pod.Namespace, restrictedPath)
			return false
		}
	}

	return true // Passes the check if no restricted API paths are accessed
}

// getAPIRestrictions loads the API restrictions from the configuration file
func getAPIRestrictions() APIRestrictions {
	// This function would typically load restrictions from a config file like security-policies.yaml
	// Here we define them statically for demonstration purposes
	return APIRestrictions{
		RestrictedAPIPaths: []string{
			"/api/v1/namespaces/*/pods/exec",
			"/api/v1/nodes/*/proxy",
		},
	}
}

// APIRestrictions defines a structure for restricted API paths
type APIRestrictions struct {
	RestrictedAPIPaths []string `yaml:"restrictedAPIPaths"`
}
