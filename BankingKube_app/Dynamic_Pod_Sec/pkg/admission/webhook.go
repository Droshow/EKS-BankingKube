package admission

import (
	"encoding/json"
	admissionv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"github.com/Droshow/EKS-BankingKube/BankingKube_app/Dynamic_Pod_Sec/pkg/admission/context_capabilities"
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
	if !checkPrivilegedContainers(request) {
		allowed = false
		result = &metav1.Status{
			Message: "Pod contains privileged containers, which is not allowed.",
		}
	}

	if !checkHostPath(request) {
		allowed = false
		result = &metav1.Status{
			Message: "Pod contains disallowed host paths.",
		}
	}

	if !checkReadOnlyRootFilesystem(request) {
		allowed = false
		result = &metav1.Status{
			Message: "Pod contains containers with writable root filesystems.",
		}
	}

	// Return the response
	return &admissionv1.AdmissionResponse{
		Allowed: allowed,
		Result:  result,
	}
}
