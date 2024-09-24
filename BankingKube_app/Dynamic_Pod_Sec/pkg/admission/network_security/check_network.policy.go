package network_security

import (
	"encoding/json"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"log"
)

// CheckNetworkPolicy ensures that there are network policies applied to the namespace or pod
func CheckNetworkPolicy(request *admissionv1.AdmissionRequest) bool {
	// Parse the Pod object from the request
	pod := &corev1.Pod{}
	err := json.Unmarshal(request.Object.Raw, pod)
	if err != nil {
		log.Println("Failed to parse Pod object:", err)
		return false // Return false if pod parsing fails
	}

	// Check if there are network policies in the namespace
	networkPolicies, err := getNetworkPolicies(pod.Namespace)
	if err != nil {
		log.Println("Failed to retrieve network policies:", err)
		return false // Return false if we cannot fetch network policies
	}

	// Check if any network policy applies to the pod
	for _, policy := range networkPolicies {
		if isNetworkPolicyApplicable(policy, pod) {
			return true // Return true if at least one network policy is applicable
		}
	}

	log.Printf("Pod %s in namespace %s has no network policies applied, which is not allowed.", pod.Name, pod.Namespace)
	return false // Return false if no network policies apply to the pod
}

// getNetworkPolicies retrieves the network policies in the given namespace
func getNetworkPolicies(namespace string) ([]unstructured.Unstructured, error) {
	// Here we would use a Kubernetes client to list network policies
	// This is a placeholder implementation
	// Replace with actual Kubernetes client logic
	return []unstructured.Unstructured{}, nil
}

// isNetworkPolicyApplicable checks if a network policy is applicable to a pod
func isNetworkPolicyApplicable(policy unstructured.Unstructured, pod *corev1.Pod) bool {
	// Check if policy applies to the pod based on podSelector
	podSelector, found, err := unstructured.NestedStringMap(policy.Object, "spec", "podSelector", "matchLabels")
	if err != nil || !found {
		log.Println("Failed to retrieve podSelector from policy")
		return false
	}

	if !matchesPodSelector(podSelector, pod) {
		return false
	}

	// Check if the policy has required ingress and egress rules
	if !hasRequiredIngressEgress(policy, pod) {
		return false
	}

	// Additional best practice checks can be added here
	return true
}

// matchesPodSelector checks if the pod labels match the given selector
func matchesPodSelector(selector map[string]string, pod *corev1.Pod) bool {
	for key, value := range selector {
		if pod.Labels[key] != value {
			return false
		}
	}
	return true
}

// hasRequiredIngressEgress ensures that the policy contains necessary ingress and egress rules
func hasRequiredIngressEgress(policy unstructured.Unstructured, pod *corev1.Pod) bool {
	// Implement logic to validate ingress and egress rules
	// For example, check if the policy enforces default deny for ingress and egress
	ingressRules, found, err := unstructured.NestedSlice(policy.Object, "spec", "ingress")
	if err != nil || !found {
		log.Println("No ingress rules found in policy")
		return false
	}

	egressRules, found, err := unstructured.NestedSlice(policy.Object, "spec", "egress")
	if err != nil || !found {
		log.Println("No egress rules found in policy")
		return false
	}

	// Ensure default deny-all is present
	if !containsDefaultDeny(ingressRules, egressRules) {
		log.Println("Default deny-all policy is missing")
		return false
	}

	return true
}

// containsDefaultDeny checks if default deny-all rules are present in the ingress and egress rules
func containsDefaultDeny(ingressRules []interface{}, egressRules []interface{}) bool {
	denyAllFound := false

	for _, rule := range ingressRules {
		// Implement logic to check for deny-all rule in ingress
		if ruleHasDenyAll(rule) {
			denyAllFound = true
			break
		}
	}

	if !denyAllFound {
		log.Println("Ingress deny-all rule not found")
		return false
	}

	denyAllFound = false

	for _, rule := range egressRules {
		// Implement logic to check for deny-all rule in egress
		if ruleHasDenyAll(rule) {
			denyAllFound = true
			break
		}
	}

	if !denyAllFound {
		log.Println("Egress deny-all rule not found")
		return false
	}

	return true
}

// ruleHasDenyAll checks if a given rule contains deny-all criteria
func ruleHasDenyAll(rule interface{}) bool {
	ruleMap, ok := rule.(map[string]interface{})
	if !ok {
		return false
	}

	ports, found := ruleMap["ports"]
	if found && len(ports.([]interface{})) == 0 {
		return true
	}

	return false
}
