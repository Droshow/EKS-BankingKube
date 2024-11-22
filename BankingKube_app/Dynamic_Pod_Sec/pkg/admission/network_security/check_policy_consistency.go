package network_security

import (
	"log"
	"net"
	"os"

	"gopkg.in/yaml.v2"
)

// Defines a structure for egress network policies
type accessEgressPolicy struct {
	AllowedEgressCIDRs []string `yaml:"allowedEgressCIDRs"`
}

// Defines a structure for ingress network policies
type accessIngressPolicy struct {
	AllowedIngressCIDRs []string `yaml:"allowedIngressCIDRs"`
}

// Defines a structure for allowed overlapping egress CIDRs
type AllowedOverlappingEgressCIDRs struct {
	CIDRs []string `yaml:"allowedOverlappingEgressCIDRs"`
}

// Defines a structure for allowed overlapping ingress CIDRs
type AllowedOverlappingIngressCIDRs struct {
	CIDRs []string `yaml:"allowedOverlappingIngressCIDRs"`
}

// ConsistencyPolicy defines the structure for consistency checks in network policies
type ConsistencyPolicy struct {
	EgressPolicy                   accessEgressPolicy             `yaml:"accessEgressPolicy"`
	IngressPolicy                  accessIngressPolicy            `yaml:"accessIngressPolicy"`
	AllowedOverlappingEgressCIDRs  AllowedOverlappingEgressCIDRs  `yaml:"allowedOverlappingEgressCIDRs"`
	AllowedOverlappingIngressCIDRs AllowedOverlappingIngressCIDRs `yaml:"allowedOverlappingIngressCIDRs"`
}

// CheckPolicyConsistency validates the consistency of security policies
func CheckPolicyConsistency() bool {
	policies, err := getConsistencyPolicy()
	if err != nil {
		log.Println("Failed to load consistency policies:", err)
		return false
	}

	if !checkCIDRConsistency(policies.EgressPolicy.AllowedEgressCIDRs, policies.IngressPolicy.AllowedIngressCIDRs, policies.AllowedOverlappingEgressCIDRs.CIDRs, policies.AllowedOverlappingIngressCIDRs.CIDRs) {
		log.Println("CIDR ranges in egress and ingress policies are inconsistent")
		return false
	}

	log.Println("Security policies are consistent")
	return true
}

// getConsistencyPolicy loads the consistency policy from the configuration file
func getConsistencyPolicy() (*ConsistencyPolicy, error) {
	configPath := os.Getenv("SECURITY_POLICIES_PATH")
	if configPath == "" {
		configPath = "configs/security-policies.yaml"
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	var policies struct {
		ConsistencyPolicy ConsistencyPolicy `yaml:"consistencyPolicy"`
	}
	err = yaml.Unmarshal(data, &policies)
	return &policies.ConsistencyPolicy, err
}

// checkCIDRConsistency checks for overlapping CIDR ranges between egress and ingress policies
func checkCIDRConsistency(egressCIDRs, ingressCIDRs, allowedEgressOverlaps, allowedIngressOverlaps []string) bool {
	for _, egressCIDR := range egressCIDRs {
		_, egressNet, err := net.ParseCIDR(egressCIDR)
		if err != nil {
			log.Printf("Invalid egress CIDR format: %s\n", egressCIDR)
			return false
		}
		for _, ingressCIDR := range ingressCIDRs {
			_, ingressNet, err := net.ParseCIDR(ingressCIDR)
			if err != nil {
				log.Printf("Invalid ingress CIDR format: %s\n", ingressCIDR)
				return false
			}
			if isAllowedOverlap(egressCIDR, ingressCIDR, allowedEgressOverlaps, allowedIngressOverlaps) {
				continue
			}
			if egressNet.Contains(ingressNet.IP) || ingressNet.Contains(egressNet.IP) {
				log.Printf("Overlapping CIDR ranges found: %s and %s\n", egressCIDR, ingressCIDR)
				return false
			}
		}
	}
	return true
}

// isAllowedOverlap checks if the overlap between egress and ingress CIDRs is allowed
func isAllowedOverlap(egressCIDR, ingressCIDR string, allowedEgressOverlaps, allowedIngressOverlaps []string) bool {
	for _, allowedEgress := range allowedEgressOverlaps {
		for _, allowedIngress := range allowedIngressOverlaps {
			if egressCIDR == allowedEgress && ingressCIDR == allowedIngress {
				return true
			}
		}
	}
	return false
}
