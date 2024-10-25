package volume_security

import (
	"encoding/json"
	"log"
	"os"

	"gopkg.in/yaml.v2"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
)

// VolumeSecurity defines a structure for volume security policies
type VolumeSecurity struct {
	DisallowedHostPaths      []string `yaml:"disallowedHostPaths"`
	RestrictedStorageClasses []string `yaml:"restrictedStorageClasses"`
}

// CheckHostPath checks if the pod has any disallowed hostPath volumes
func CheckHostPath(request *admissionv1.AdmissionRequest) bool {
	// Parse the Pod object from the request
	pod := &corev1.Pod{}
	err := json.Unmarshal(request.Object.Raw, pod)
	if err != nil {
		log.Println("Failed to parse pod object:", err)
		return false // Fails the validation if the pod can't be parsed
	}

	// Retrieve the volume security policies
	volumeSecurity, err := getVolumeSecurity()
	if err != nil {
		log.Println("Failed to load volume security policies:", err)
		return false
	}

	// Check for disallowed hostPath volumes
	for _, volume := range pod.Spec.Volumes {
		if volume.HostPath != nil {
			for _, disallowedPath := range volumeSecurity.DisallowedHostPaths {
				if volume.HostPath.Path == disallowedPath {
					log.Println("Disallowed hostPath volume found:", volume.HostPath.Path)
					return false // Return false if any disallowed hostPath volume is found
				}
			}
		}
	}

	// Passes the check if no disallowed hostPath volumes are found
	return true
}

// getVolumeSecurity loads the volume security policies from the configuration file
func getVolumeSecurity() (*VolumeSecurity, error) {
	configPath := os.Getenv("SECURITY_POLICIES_PATH")
	if configPath == "" {
		configPath = "configs/security-policies.yaml" // Default path
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var policies struct {
		VolumeSecurity VolumeSecurity `yaml:"volumeSecurity"`
	}

	err = yaml.Unmarshal(data, &policies)
	if err != nil {
		return nil, err
	}

	return &policies.VolumeSecurity, nil
}
