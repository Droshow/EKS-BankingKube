package admission

import (
	"encoding/json"
	"net/http"

	"github.com/Droshow/EKS-BankingKube/BankingKube_app/Dynamic_Pod_Sec/pkg/admission/api_restrictions"
	"github.com/Droshow/EKS-BankingKube/BankingKube_app/Dynamic_Pod_Sec/pkg/admission/context_capabilities"
	"github.com/Droshow/EKS-BankingKube/BankingKube_app/Dynamic_Pod_Sec/pkg/admission/image_security"
	"github.com/Droshow/EKS-BankingKube/BankingKube_app/Dynamic_Pod_Sec/pkg/admission/network_security"
	"github.com/Droshow/EKS-BankingKube/BankingKube_app/Dynamic_Pod_Sec/pkg/admission/resource_limits"
	"github.com/Droshow/EKS-BankingKube/BankingKube_app/Dynamic_Pod_Sec/pkg/admission/volume_security"
	"github.com/Droshow/EKS-BankingKube/BankingKube_app/Dynamic_Pod_Sec/pkg/admission/rbac_checks"

	admissionv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// HandleAdmissionRequest handles incoming admission requests
func HandleAdmissionRequest(w http.ResponseWriter, r *http.Request) {
	var admissionReview admissionv1.AdmissionReview
	err := json.NewDecoder(r.Body).Decode(&admissionReview)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Main validation logic that calls individual security checks
	response := validatePod(admissionReview.Request)

	admissionReview.Response = response
	admissionReview.Response.UID = admissionReview.Request.UID

	resp, err := json.Marshal(admissionReview)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

// validatePod is the main function that calls individual validation functions
func validatePod(request *admissionv1.AdmissionRequest) *admissionv1.AdmissionResponse {
	allowed := true
	result := &metav1.Status{
		Message: "Pod validation passed",
	}

	// Check for Pod Context Security and Capabilities
	if !context_capabilities.CheckPodSecurityContext(request) {
		allowed = false
		result = &metav1.Status{
			Message: "Pod contains privileged containers, which is not allowed.",
		}
	}

	if !context_capabilities.CheckCapabilities(request) {
		allowed = false
		result = &metav1.Status{
			Message: "Pod contains containers with disallowed capabilities.",
		}
	}
	// Check for Pod Volume Security
	if !volume_security.CheckHostPath(request) {
		allowed = false
		result = &metav1.Status{
			Message: "Pod contains disallowed host paths.",
		}
	}
	// Check for Pod Network Policy Security
	if !network_security.CheckPolicyConsistency() {
		allowed = false
		result = &metav1.Status{
			Message: "Security policies are inconsistent.",
		}
	}

	if !network_security.CheckNetworkPolicy(request) { // Integrate CheckNetworkPolicy
		allowed = false
		result = &metav1.Status{
			Message: "Pod does not comply with network policies.",
		}
	}

	if !network_security.CheckHostNetwork(request) {
		allowed = false
		result = &metav1.Status{
			Message: "Pod is using the host network, which is disallowed.",
		}
	}
	if !network_security.CheckEgress(request) {
		allowed = false
		result = &metav1.Status{
			Message: "Pod  has an egress route that violates policy.",
		}
	}
	if !network_security.CheckIngress(request) {
		allowed = false
		result = &metav1.Status{
			Message: "Pod  has an ingress route that violates policy.",
		}
	}

	// Check for API Restrictions & Service Accounts
	if !api_restrictions.CheckAPIAccess(request) { // Integrate CheckAPIAccess
		allowed = false
		result = &metav1.Status{
			Message: "Pod is attempting to access restricted API paths.",
		}
	}

	if !api_restrictions.CheckServiceAccount(request) { // Integrate CheckServiceAccount
		allowed = false
		result = &metav1.Status{
			Message: "Pod is using a restricted service account.",
		}
	}
	// Check for Pod Image Security
	if !image_security.CheckImageRegistry(request) {
		allowed = false
		result = &metav1.Status{
			Message: "Pod is using an image from a disallowed registry.",
		}
	}

	if !image_security.CheckImageSigning(request) {
		allowed = false
		result = &metav1.Status{
			Message: "Pod is using an unsigned image.",
		}
	}

	if !image_security.CheckImageTags(request) {
		allowed = false
		result = &metav1.Status{
			Message: "Pod is using an image with a disallowed tag.",
		}
	}
	// Check for RBAC Policies
	if !rbac_checks.CheckRBACBinding(request) { 
		allowed = false
		result = &metav1.Status{
			Message: "ClusterRoleBinding uses restricted ClusterRole.",
		}
	}


	// Check for Pod Resource Limits
	if !resource_limits.CheckResourceLimits(request) { // Integrate CheckResourceLimits
		allowed = false
		result = &metav1.Status{
			Message: "Pod contains containers without resource limits.",
		}
	}

	if !resource_limits.CheckResourceRequests(request) { // Integrate CheckResourceRequests
		allowed = false
		result = &metav1.Status{
			Message: "Pod contains containers without resource requests.",
		}
	}

	// Return the response
	return &admissionv1.AdmissionResponse{
		Allowed: allowed,
		Result:  result,
	}
}
