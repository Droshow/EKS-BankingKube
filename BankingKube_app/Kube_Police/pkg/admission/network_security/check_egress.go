package network_security

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	"log"
	"os"
)

// EgressPolicy defines a structure for egress network policies
type EgressPolicy struct {
	AllowedEgressCIDRs []string `yaml:"allowedEgressCIDRs"`
}

// CheckEgress validates if a pod has correct egress restrictions
func CheckEgress(request *admissionv1.AdmissionRequest) bool {
	pod := &corev1.Pod{}
	err := json.Unmarshal(request.Object.Raw, pod)
	if err != nil {
		log.Println("Failed to parse pod object:", err)
		return false
	}

	egressPolicy, err := getEgressPolicy()
	if err != nil {
		log.Println("Failed to load egress policy:", err)
		return false
	}

	for _, egressCIDR := range egressPolicy.AllowedEgressCIDRs {
		if !isEgressAllowed(pod, egressCIDR) {
			log.Printf("Pod %s in namespace %s has an egress route that violates policy\n", pod.Name, pod.Namespace)
			return false
		}
	}
	return true
}

func getEgressPolicy() (*EgressPolicy, error) {
	configPath := os.Getenv("SECURITY_POLICIES_PATH")
	if configPath == "" {
		configPath = "configs/security-policies.yaml"
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	var policies struct {
		EgressPolicy EgressPolicy `yaml:"egressPolicy"`
	}
	err = yaml.Unmarshal(data, &policies)
	return &policies.EgressPolicy, err
}

func isEgressAllowed(pod *corev1.Pod, egressCIDR string) bool {
	// Implement logic to validate if egress is within allowed CIDR ranges
	return true
}
