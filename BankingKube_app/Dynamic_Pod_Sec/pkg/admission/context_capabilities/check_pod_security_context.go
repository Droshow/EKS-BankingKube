package context_capabilities

import (
	"encoding/json"
	"log"
	"os"

	"gopkg.in/yaml.v2"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
)

// PodSecurityContext defines a structure for pod security context policies
type PodSecurityContext struct {
	AllowPrivilegeEscalation bool `yaml:"allowPrivilegeEscalation"`
	RunAsNonRoot             bool `yaml:"runAsNonRoot"`
	ReadOnlyRootFilesystem   bool `yaml:"readOnlyRootFilesystem"`
}

// CheckPodSecurityContext checks if the pod complies with the pod security context policies
func CheckPodSecurityContext(request *admissionv1.AdmissionRequest) bool {
	// Parse the Pod object from the request
	pod := &corev1.Pod{}
	err := json.Unmarshal(request.Object.Raw, pod)
	if err != nil {
		log.Println("Failed to parse pod object:", err)
		return false // Fails the validation if the pod can't be parsed
	}

	// Retrieve the security policies for pod security context
	podSecurityContext, err := getPodSecurityContext()
	if err != nil {
		log.Println("Failed to load pod security context policies:", err)
		return false
	}

	// Check for privileged containers
	for _, container := range pod.Spec.Containers {
		if container.SecurityContext == nil {
			log.Println("Warning: Container", container.Name, "does not have a SecurityContext defined")
			return false // Consider it a security violation if SecurityContext is not defined
		}
		if container.SecurityContext.Privileged != nil && *container.SecurityContext.Privileged {
			return false // Return false if any container is privileged
		}
		if container.SecurityContext.AllowPrivilegeEscalation != nil && *container.SecurityContext.AllowPrivilegeEscalation != podSecurityContext.AllowPrivilegeEscalation {
			log.Println("Container", container.Name, "has AllowPrivilegeEscalation set to", *container.SecurityContext.AllowPrivilegeEscalation, "which does not match the policy")
			return false // Return false if AllowPrivilegeEscalation does not match the policy
		}
		if container.SecurityContext.RunAsNonRoot != nil && *container.SecurityContext.RunAsNonRoot != podSecurityContext.RunAsNonRoot {
			log.Println("Container", container.Name, "has RunAsNonRoot set to", *container.SecurityContext.RunAsNonRoot, "which does not match the policy")
			return false // Return false if RunAsNonRoot does not match the policy
		}
		if container.SecurityContext.ReadOnlyRootFilesystem != nil && *container.SecurityContext.ReadOnlyRootFilesystem != podSecurityContext.ReadOnlyRootFilesystem {
			log.Println("Container", container.Name, "has ReadOnlyRootFilesystem set to", *container.SecurityContext.ReadOnlyRootFilesystem, "which does not match the policy")
			return false // Return false if ReadOnlyRootFilesystem does not match the policy
		}
	}

	// Check for privileged init containers
	for _, initContainer := range pod.Spec.InitContainers {
		if initContainer.SecurityContext == nil {
			log.Println("Warning: Init container", initContainer.Name, "does not have a SecurityContext defined")
			return false // Consider it a security violation if SecurityContext is not defined
		}
		if initContainer.SecurityContext.Privileged != nil && *initContainer.SecurityContext.Privileged {
			return false // Return false if any init container is privileged
		}
		if initContainer.SecurityContext.AllowPrivilegeEscalation != nil && *initContainer.SecurityContext.AllowPrivilegeEscalation != podSecurityContext.AllowPrivilegeEscalation {
			log.Println("Init container", initContainer.Name, "has AllowPrivilegeEscalation set to", *initContainer.SecurityContext.AllowPrivilegeEscalation, "which does not match the policy")
			return false // Return false if AllowPrivilegeEscalation does not match the policy
		}
		if initContainer.SecurityContext.RunAsNonRoot != nil && *initContainer.SecurityContext.RunAsNonRoot != podSecurityContext.RunAsNonRoot {
			log.Println("Init container", initContainer.Name, "has RunAsNonRoot set to", *initContainer.SecurityContext.RunAsNonRoot, "which does not match the policy")
			return false // Return false if RunAsNonRoot does not match the policy
		}
		if initContainer.SecurityContext.ReadOnlyRootFilesystem != nil && *initContainer.SecurityContext.ReadOnlyRootFilesystem != podSecurityContext.ReadOnlyRootFilesystem {
			log.Println("Init container", initContainer.Name, "has ReadOnlyRootFilesystem set to", *initContainer.SecurityContext.ReadOnlyRootFilesystem, "which does not match the policy")
			return false // Return false if ReadOnlyRootFilesystem does not match the policy
		}
	}

	// Passes the check if no privileged containers or init containers are found
	return true
}

// getPodSecurityContext loads the pod security context policies from the configuration file
func getPodSecurityContext() (*PodSecurityContext, error) {
	configPath := os.Getenv("SECURITY_POLICIES_PATH")
	if configPath == "" {
		configPath = "configs/security-policies.yaml" // Default path
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var policies struct {
		PodSecurityContext PodSecurityContext `yaml:"podSecurityContext"`
	}

	err = yaml.Unmarshal(data, &policies)
	if err != nil {
		return nil, err
	}

	return &policies.PodSecurityContext, nil
}
