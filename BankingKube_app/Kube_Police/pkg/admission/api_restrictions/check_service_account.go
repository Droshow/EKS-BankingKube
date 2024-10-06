package api_restrictions

import (
	"encoding/json"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	"log"
)

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
	serviceAccountRestrictions := getServiceAccountRestrictions()

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
func getServiceAccountRestrictions() ServiceAccountRestrictions {
	// This function would typically load restrictions from a config file like security-policies.yaml
	// Here we define them statically for demonstration purposes
	return ServiceAccountRestrictions{
		RestrictedServiceAccounts: []string{
			"default",
			"admin",
		},
	}
}

// ServiceAccountRestrictions defines a structure for restricted service accounts
type ServiceAccountRestrictions struct {
	RestrictedServiceAccounts []string `yaml:"restrictedServiceAccounts"`
}
