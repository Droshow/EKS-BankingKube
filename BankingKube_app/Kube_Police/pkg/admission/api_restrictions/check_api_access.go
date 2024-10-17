package api_restrictions

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
)

// APIRestrictions defines a structure for restricted API paths and service accounts
type APIRestrictions struct {
	RestrictedServiceAccounts []string `yaml:"restrictedServiceAccounts"`
	RestrictedAPIPaths        []string `yaml:"restrictedAPIPaths"`
}

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
	apiRestrictions, err := getAPIRestrictions()
	if err != nil {
		log.Println("Failed to load API restrictions:", err)
		return false
	}

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
func getAPIRestrictions() (*APIRestrictions, error) {
	configPath := os.Getenv("SECURITY_POLICIES_PATH")
	if configPath == "" {
		configPath = "configs/security-policies.yaml" // Default path
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var policies struct {
		APIRestrictions APIRestrictions `yaml:"apiRestrictions"`
	}

	err = yaml.Unmarshal(data, &policies)
	if err != nil {
		return nil, err
	}

	return &policies.APIRestrictions, nil
}
