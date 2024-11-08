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

// IngressPolicy defines a structure for ingress network policies
type IngressPolicy struct {
	AllowedIngressCIDRs []string `yaml:"allowedIngressCIDRs"`
}

// CheckIngress validates if a pod has correct ingress restrictions
func CheckIngress(request *admissionv1.AdmissionRequest) bool {
	pod := &corev1.Pod{}
	err := json.Unmarshal(request.Object.Raw, pod)
	if err != nil {
		log.Println("Failed to parse pod object:", err)
		return false
	}

	ingressPolicy, err := getIngressPolicy()
	if err != nil {
		log.Println("Failed to load ingress policy:", err)
		return false
	}

	for _, ingressCIDR := range ingressPolicy.AllowedIngressCIDRs {
		if !isIngressAllowed(pod, ingressCIDR) {
			log.Printf("Pod %s in namespace %s has an ingress route that violates policy\n", pod.Name, pod.Namespace)
			return false
		}
	}
	return true
}

func getIngressPolicy() (*IngressPolicy, error) {
	configPath := os.Getenv("SECURITY_POLICIES_PATH")
	if configPath == "" {
		configPath = "configs/security-policies.yaml"
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	var policies struct {
		IngressPolicy IngressPolicy `yaml:"ingressPolicy"`
	}
	err = yaml.Unmarshal(data, &policies)
	return &policies.IngressPolicy, err
}

// isIngressAllowed checks if the pod's ingress traffic is within the allowed CIDR ranges
func isIngressAllowed(pod *corev1.Pod, ingressCIDR string) bool {
	// Parse the allowed ingress CIDR
	_, allowedCIDR, err := net.ParseCIDR(ingressCIDR)
	if err != nil {
		log.Printf("Invalid CIDR format: %s\n", ingressCIDR)
		return false
	}

	// Check pod annotations for ingress IPs (assuming annotations are used to specify ingress IPs)
	ingressIPs, ok := pod.Annotations["ingressIPs"]
	if !ok {
		log.Printf("Pod %s in namespace %s does not have ingress IPs specified\n", pod.Name, pod.Namespace)
		return false
	}

	// Split the ingress IPs and check each one
	for _, ip := range strings.Split(ingressIPs, ",") {
		parsedIP := net.ParseIP(ip)
		if parsedIP == nil {
			log.Printf("Invalid IP format: %s\n", ip)
			return false
		}

		// Check if the IP is within the allowed CIDR range
		if !allowedCIDR.Contains(parsedIP) {
			log.Printf("IP %s is not within the allowed CIDR range %s\n", ip, ingressCIDR)
			return false
		}
	}

	return true
}
