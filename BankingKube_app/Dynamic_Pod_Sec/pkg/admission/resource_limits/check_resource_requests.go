package resource_limits

import (
    "encoding/json"
    "log"
    "os"

    "gopkg.in/yaml.v2"
    admissionv1 "k8s.io/api/admission/v1"
    corev1 "k8s.io/api/core/v1"
)

// EnforceResourceRequests defines a structure for the enforceResourceRequests policy
type EnforceResourceRequests struct {
    EnforceResourceRequests bool `yaml:"enforceResourceRequests"`
}

// CheckResourceRequests validates if a pod's containers have resource requests defined
func CheckResourceRequests(request *admissionv1.AdmissionRequest) bool {
    // Parse the Pod object from the request
    pod := &corev1.Pod{}
    err := json.Unmarshal(request.Object.Raw, pod)
    if err != nil {
        log.Println("Failed to parse pod object:", err)
        return false // Fails the validation if the pod can't be parsed
    }

    // Retrieve the enforceResourceRequests policy
    enforceResourceRequests, err := getEnforceResourceRequests()
    if err != nil {
        log.Println("Failed to load enforceResourceRequests policy:", err)
        return false
    }

    // Check if resource requests enforcement is enabled
    if !enforceResourceRequests.EnforceResourceRequests {
        return true // Passes the check if enforcement is not enabled
    }

    // Check if the pod's containers have resource requests defined
    for _, container := range pod.Spec.Containers {
        if container.Resources.Requests == nil {
            log.Printf("Pod %s in namespace %s has a container without resource requests: %s\n", pod.Name, pod.Namespace, container.Name)
            return false
        }
    }

    return true // Passes the check if all containers have resource requests defined
}

// getEnforceResourceRequests loads the enforceResourceRequests policy from the configuration file
func getEnforceResourceRequests() (*EnforceResourceRequests, error) {
    configPath := os.Getenv("SECURITY_POLICIES_PATH")
    if configPath == "" {
        configPath = "configs/security-policies.yaml" // Default path
    }

    data, err := os.ReadFile(configPath)
    if err != nil {
        return nil, err
    }

    var policies struct {
        EnforceResourceRequests EnforceResourceRequests `yaml:"resourceLimits"`
    }

    err = yaml.Unmarshal(data, &policies)
    if err != nil {
        return nil, err
    }

    return &policies.EnforceResourceRequests, nil
}