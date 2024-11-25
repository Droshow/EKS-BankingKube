package image_security

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"

	"gopkg.in/yaml.v2"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
)

// RequireImageSigning defines a structure for the requireImageSigning policy
type RequireImageSigning struct {
	RequireImageSigning bool `yaml:"requireImageSigning"`
}

// CheckImageSigning validates if a pod's images are signed and verified
func CheckImageSigning(request *admissionv1.AdmissionRequest) bool {
	// Parse the Pod object from the request
	pod := &corev1.Pod{}
	err := json.Unmarshal(request.Object.Raw, pod)
	if err != nil {
		log.Println("Failed to parse pod object:", err)
		return false // Fails the validation if the pod can't be parsed
	}

	// Retrieve the requireImageSigning policy
	requireImageSigning, err := getRequireImageSigning()
	if err != nil {
		log.Println("Failed to load requireImageSigning policy:", err)
		return false
	}

	// Check if image signing enforcement is enabled
	if !requireImageSigning.RequireImageSigning {
		return true // Passes the check if enforcement is not enabled
	}

	// Check if the pod's containers are using signed images
	for _, container := range pod.Spec.Containers {
		if !isImageSigned(container.Image) {
			log.Printf("Pod %s in namespace %s is using an unsigned image: %s\n", pod.Name, pod.Namespace, container.Image)
			return false
		}
	}

	return true // Passes the check if all images are signed
}

// getRequireImageSigning loads the requireImageSigning policy from the configuration file
func getRequireImageSigning() (*RequireImageSigning, error) {
	configPath := os.Getenv("SECURITY_POLICIES_PATH")
	if configPath == "" {
		configPath = "configs/security-policies.yaml" // Default path
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var policies struct {
		ImageSecurity struct {
			RequireImageSigning RequireImageSigning `yaml:"requireImageSigning"`
		} `yaml:"imageSecurity"`
	}

	err = yaml.Unmarshal(data, &policies)
	if err != nil {
		return nil, err
	}

	return &policies.ImageSecurity.RequireImageSigning, nil
}

// isImageSigned checks if the image is signed using cosign
func isImageSigned(image string) bool {
	// Get the public key path from the environment variable
	publicKeyPath := os.Getenv("COSIGN_PUBLIC_KEY_PATH")
	if publicKeyPath == "" {
		log.Println("COSIGN_PUBLIC_KEY_PATH environment variable is not set")
		return false
	}

	// Construct the cosign verify command
	cmd := exec.Command("cosign", "verify", "--key", publicKeyPath, image)

	// Run the command and capture the output
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Failed to verify image signature for %s: %v\nOutput: %s", image, err, string(output))
		return false
	}

	log.Printf("Successfully verified image signature for %s\n", image)
	return true
}
