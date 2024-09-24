package admission

import (
	"encoding/json"
	"net/http"

	"github.com/Droshow/EKS-BankingKube/BankingKube_app/Dynamic_Pod_Sec/pkg/admission/context_capabilities"
	"github.com/Droshow/EKS-BankingKube/BankingKube_app/Dynamic_Pod_Sec/pkg/admission/volume_security"
	"github.com/Droshow/EKS-BankingKube/BankingKube_app/Dynamic_Pod_Sec/pkg/admission/network_security"
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

	// Run individual validation checks
	if !context_capabilities.CheckPrivilegedContainers(request) {
		allowed = false
		result = &metav1.Status{
			Message: "Pod contains privileged containers, which is not allowed.",
		}
	}

	if !context_capabilities.CheckReadOnlyRoot(request) {
		allowed = false
		result = &metav1.Status{
			Message: "Pod contains containers with writable root filesystems.",
		}
	}

	if !context_capabilities.CheckRunAsUser(request) {
		allowed = false
		result = &metav1.Status{
			Message: "Pod does not meet user/group security context requirements.",
		}
	}

	if !context_capabilities.CheckCapabilities(request) {
		allowed = false
		result = &metav1.Status{
			Message: "Pod contains containers with disallowed capabilities.",
		}
	}

	if !context_capabilities.ValidateCapability(request) {
		allowed = false
		result = &metav1.Status{
			Message: "Pod contains containers with unvalidated capabilities.",
		}
	}

	if !volume_security.CheckHostPath(request) {
		allowed = false
		result = &metav1.Status{
			Message: "Pod contains disallowed host paths.",
		}
	}

	if !network_security.CheckNetworkPolicy(request) { // Integrate CheckNetworkPolicy
        allowed = false
        result = &metav1.Status{
            Message: "Pod does not comply with network policies.",
        }
    }

	// Return the response
	return &admissionv1.AdmissionResponse{
		Allowed: allowed,
		Result:  result,
	}
}
