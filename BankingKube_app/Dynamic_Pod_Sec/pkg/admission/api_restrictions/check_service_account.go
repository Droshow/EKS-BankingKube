package api_restrictions

import (
	"encoding/json"
	"log"
	"os"

	"gopkg.in/yaml.v2"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
)

// ServiceAccountRestrictions defines a structure for restricted service accounts
type ServiceAccountRestrictions struct {
	RestrictedServiceAccounts []string `yaml:"restrictedServiceAccounts"`
}

// CheckServiceAccount ensures that pods do not use restricted service accounts
func CheckServiceAccount(request *admissionv1.AdmissionRequest) bool {
	// Parse the Pod object from the request
	pod := &corev1.Pod{}
	err := json.Unmarshal(request.Object.Raw, pod)
	if err != nil {
		log.Println("Failed to parse pod object:", err)
		return false // Fails the validation if the pod can't be parsed
	}

	// Retrieve the security policies for service account usage
	serviceAccountRestrictions, err := getServiceAccountRestrictions()
	if err != nil {
		log.Println("Failed to load service account restrictions:", err)
		return false
	}

	// Check if the pod is using a restricted service account
	for _, restrictedAccount := range serviceAccountRestrictions.RestrictedServiceAccounts {
		if pod.Spec.ServiceAccountName == restrictedAccount {
			log.Printf("Pod %s in namespace %s is using a restricted service account: %s\n", pod.Name, pod.Namespace, restrictedAccount)
			return false
		}
	}

	return true // Passes the check if no restricted service accounts are used
}

// getServiceAccountRestrictions loads the service account restrictions from the configuration file
func getServiceAccountRestrictions() (*ServiceAccountRestrictions, error) {
	configPath := os.Getenv("SECURITY_POLICIES_PATH")
	if configPath == "" {
		configPath = "configs/security-policies.yaml" // Default path
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var policies struct {
		ServiceAccountRestrictions ServiceAccountRestrictions `yaml:"serviceAccountRestrictions"`
	}

	err = yaml.Unmarshal(data, &policies)
	if err != nil {
		return nil, err
	}

	return &policies.ServiceAccountRestrictions, nil
}
