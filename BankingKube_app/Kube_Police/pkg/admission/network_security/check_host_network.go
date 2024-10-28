package network_security

import (
	"encoding/json"
	"log"
	"os"

	"gopkg.in/yaml.v2"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
)

// HostNetworkPolicy defines a structure for host network policies
type HostNetworkPolicy struct {
	AllowHostNetwork bool `yaml:"allowHostNetwork"`
}

// CheckHostNetwork validates if a pod is using the host network
func CheckHostNetwork(request *admissionv1.AdmissionRequest) bool {
	// Parse the Pod object from the request
	pod := &corev1.Pod{}
	err := json.Unmarshal(request.Object.Raw, pod)
	if err != nil {
		log.Println("Failed to parse pod object:", err)
		return false // Fails the validation if the pod can't be parsed
	}

	// Retrieve the host network policy
	hostNetworkPolicy, err := getHostNetworkPolicy()
	if err != nil {
		log.Println("Failed to load host network policy:", err)
		return false
	}

	// Check if the pod is using the host network
	if pod.Spec.HostNetwork && !hostNetworkPolicy.AllowHostNetwork {
		log.Printf("Pod %s in namespace %s is using the host network, which is disallowed\n", pod.Name, pod.Namespace)
		return false
	}

	return true // Passes the check if the pod is not using the host network or if it is allowed
}

// getHostNetworkPolicy loads the host network policy from the configuration file
func getHostNetworkPolicy() (*HostNetworkPolicy, error) {
	configPath := os.Getenv("SECURITY_POLICIES_PATH")
	if configPath == "" {
		configPath = "configs/security-policies.yaml" // Default path
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var policies struct {
		HostNetworkPolicy HostNetworkPolicy `yaml:"hostNetworkPolicy"`
	}

	err = yaml.Unmarshal(data, &policies)
	if err != nil {
		return nil, err
	}

	return &policies.HostNetworkPolicy, nil
}
