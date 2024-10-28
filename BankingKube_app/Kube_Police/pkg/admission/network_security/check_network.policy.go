package network_security

import (
	"encoding/json"
	"log"
	"os"

	"gopkg.in/yaml.v2"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
)

// NetworkPolicy defines a structure for network policies
type NetworkPolicy struct {
	RequiredNetworkPolicies []string `yaml:"requiredNetworkPolicies"`
}

// CheckNetworkPolicy validates if a pod is using the required network policies
func CheckNetworkPolicy(request *admissionv1.AdmissionRequest) bool {
	// Parse the Pod object from the request
	pod := &corev1.Pod{}
	err := json.Unmarshal(request.Object.Raw, pod)
	if err != nil {
		log.Println("Failed to parse pod object:", err)
		return false // Fails the validation if the pod can't be parsed
	}

	// Retrieve the network policies
	networkPolicy, err := getNetworkPolicy()
	if err != nil {
		log.Println("Failed to load network policies:", err)
		return false
	}

	// Check if the pod's network policies are in the required list
	for _, policy := range networkPolicy.RequiredNetworkPolicies {
		if !isNetworkPolicyPresent(pod, policy) {
			log.Printf("Pod %s in namespace %s is missing required network policy: %s\n", pod.Name, pod.Namespace, policy)
			return false
		}
	}

	return true // Passes the check if all required network policies are present
}

// getNetworkPolicy loads the network policies from the configuration file
func getNetworkPolicy() (*NetworkPolicy, error) {
	configPath := os.Getenv("SECURITY_POLICIES_PATH")
	if configPath == "" {
		configPath = "configs/security-policies.yaml" // Default path
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var policies struct {
		NetworkPolicy NetworkPolicy `yaml:"networkPolicy"`
	}

	err = yaml.Unmarshal(data, &policies)
	if err != nil {
		return nil, err
	}

	return &policies.NetworkPolicy, nil
}

// isNetworkPolicyPresent checks if a network policy is present in the pod's annotations
func isNetworkPolicyPresent(pod *corev1.Pod, policy string) bool {
	// Check if the pod has annotations
	if pod.Annotations == nil {
		return false
	}

	// Check if the required network policy is present in the pod's annotations
	for key, value := range pod.Annotations {
		if key == "k8s.v1.cni.cncf.io/networks" && value == policy {
			return true
		}
	}

	return false
}
