package network_security

import (
	"encoding/json"
	"log"
	"net"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
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
	// Parse the allowed egress CIDR
	_, allowedCIDR, err := net.ParseCIDR(egressCIDR)
	if err != nil {
		log.Printf("Invalid CIDR format: %s\n", egressCIDR)
		return false
	}

	// Check pod annotations for egress IPs (assuming annotations are used to specify egress IPs)
	egressIPs, ok := pod.Annotations["egressIPs"]
	if !ok {
		log.Printf("Pod %s in namespace %s does not have egress IPs specified\n", pod.Name, pod.Namespace)
		return false
	}

	// Split the egress IPs and check each one
	for _, ip := range strings.Split(egressIPs, ",") {
		parsedIP := net.ParseIP(ip)
		if parsedIP == nil {
			log.Printf("Invalid IP format: %s\n", ip)
			return false
		}

		// Check if the IP is within the allowed CIDR range
		if !allowedCIDR.Contains(parsedIP) {
			log.Printf("IP %s is not within the allowed CIDR range %s\n", ip, egressCIDR)
			return false
		}
	}

	return true
}
