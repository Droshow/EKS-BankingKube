package resource_limits

import (
	"encoding/json"
	"log"
	"os"

	"gopkg.in/yaml.v2"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	resource "k8s.io/apimachinery/pkg/api/resource"
)

// ResourceLimits defines a structure for resource limits policies
type ResourceLimits struct {
	CPULimits struct {
		Max string `yaml:"max"`
		Min string `yaml:"min"`
	} `yaml:"cpuLimits"`
	MemoryLimits struct {
		Max string `yaml:"max"`
		Min string `yaml:"min"`
	} `yaml:"memoryLimits"`
	EnforceResourceLimits bool `yaml:"enforceResourceLimits"`
}

// CheckResourceLimits validates if a pod's containers have resource limits defined and within the specified range
func CheckResourceLimits(request *admissionv1.AdmissionRequest) bool {
	// Parse the Pod object from the request
	pod := &corev1.Pod{}
	err := json.Unmarshal(request.Object.Raw, pod)
	if err != nil {
		log.Println("Failed to parse pod object:", err)
		return false // Fails the validation if the pod can't be parsed
	}

	// Retrieve the resource limits policies
	resourceLimits, err := getResourceLimits()
	if err != nil {
		log.Println("Failed to load resource limits policies:", err)
		return false
	}

	// Check if resource limits enforcement is enabled
	if !resourceLimits.EnforceResourceLimits {
		return true // Passes the check if enforcement is not enabled
	}

	// Check if the pod's containers have resource limits defined and within the specified range
	for _, container := range pod.Spec.Containers {
		if container.Resources.Limits == nil {
			log.Printf("Pod %s in namespace %s has a container without resource limits: %s\n", pod.Name, pod.Namespace, container.Name)
			return false
		}

		cpuLimit := container.Resources.Limits[corev1.ResourceCPU]
		memoryLimit := container.Resources.Limits[corev1.ResourceMemory]

		if !isWithinRange(cpuLimit, resourceLimits.CPULimits.Min, resourceLimits.CPULimits.Max) {
			log.Printf("Pod %s in namespace %s has a container with CPU limit out of range: %s\n", pod.Name, pod.Namespace, container.Name)
			return false
		}

		if !isWithinRange(memoryLimit, resourceLimits.MemoryLimits.Min, resourceLimits.MemoryLimits.Max) {
			log.Printf("Pod %s in namespace %s has a container with Memory limit out of range: %s\n", pod.Name, pod.Namespace, container.Name)
			return false
		}
	}

	return true // Passes the check if all containers have resource limits defined and within the specified range
}

// isWithinRange checks if a resource quantity is within the specified range
func isWithinRange(quantity resource.Quantity, min string, max string) bool {
	minQuantity := resource.MustParse(min)
	maxQuantity := resource.MustParse(max)

	return quantity.Cmp(minQuantity) >= 0 && quantity.Cmp(maxQuantity) <= 0
}

// getResourceLimits loads the resource limits policies from the configuration file
func getResourceLimits() (*ResourceLimits, error) {
	configPath := os.Getenv("SECURITY_POLICIES_PATH")
	if configPath == "" {
		configPath = "configs/security-policies.yaml" // Default path
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var policies struct {
		ResourceLimits ResourceLimits `yaml:"resourceLimits"`
	}

	err = yaml.Unmarshal(data, &policies)
	if err != nil {
		return nil, err
	}

	return &policies.ResourceLimits, nil
}
