package context_capabilities

import (
	"encoding/json"
	"log"
	"os"

	"gopkg.in/yaml.v2"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
)

// Capabilities defines a structure for capabilities policies
type Capabilities struct {
	AllowedCapabilities    []string `yaml:"allowedCapabilities"`
	DisallowedCapabilities []string `yaml:"disallowedCapabilities"`
	RequiredDrops          []string `yaml:"requiredDrops"`
}

// CheckCapabilities ensures that the pod does not use dangerous Linux capabilities
func CheckCapabilities(request *admissionv1.AdmissionRequest) bool {
	// Parse the Pod object from the request
	pod := &corev1.Pod{}
	err := json.Unmarshal(request.Object.Raw, pod)
	if err != nil {
		log.Println("Failed to parse pod object:", err)
		return false // Fails the validation if the pod can't be parsed
	}

	// Retrieve the capabilities policies
	capabilities, err := getCapabilities()
	if err != nil {
		log.Println("Failed to load capabilities policies:", err)
		return false
	}

	// Check capabilities for each container
	for _, container := range pod.Spec.Containers {
		if container.SecurityContext != nil && container.SecurityContext.Capabilities != nil {
			for _, cap := range container.SecurityContext.Capabilities.Add {
				if isDisallowedCapability(string(cap), capabilities.DisallowedCapabilities) {
					log.Println("Disallowed capability found in container:", container.Name, "Capability:", cap)
					return false // Return false if any disallowed capability is found
				}
			}
			if !hasDroppedAllRequiredCapabilities(container.SecurityContext.Capabilities.Drop, capabilities.RequiredDrops) {
				log.Println("Necessary capabilities not dropped in container:", container.Name)
				return false // Return false if necessary capabilities are not dropped
			}
		}
	}

	// Passes the check if no disallowed capabilities are found and all required capabilities are dropped
	return true
}

// isDisallowedCapability checks if a capability is in the list of disallowed capabilities
func isDisallowedCapability(cap string, disallowedCapabilities []string) bool {
	for _, disallowedCap := range disallowedCapabilities {
		if cap == disallowedCap {
			return true
		}
	}
	return false
}

// hasDroppedAllRequiredCapabilities checks if all required capabilities are dropped
func hasDroppedAllRequiredCapabilities(dropped []corev1.Capability, requiredDrops []string) bool {
	requiredDropsMap := make(map[string]bool)
	for _, cap := range requiredDrops {
		requiredDropsMap[cap] = true
	}

	for _, cap := range dropped {
		delete(requiredDropsMap, string(cap))
	}

	return len(requiredDropsMap) == 0
}

// getCapabilities loads the capabilities policies from the configuration file
func getCapabilities() (*Capabilities, error) {
	configPath := os.Getenv("SECURITY_POLICIES_PATH")
	if configPath == "" {
		configPath = "configs/security-policies.yaml" // Default path
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var policies struct {
		Capabilities Capabilities `yaml:"capabilities"`
	}

	err = yaml.Unmarshal(data, &policies)
	if err != nil {
		return nil, err
	}

	return &policies.Capabilities, nil
}
