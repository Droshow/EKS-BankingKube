package admission

import (
	"encoding/json"
	"net/http"

	admissionv1 "k8s.io/api/admission/v1"
)

// HandleAdmissionRequest handles incoming admission requests
func HandleAdmissionRequest(w http.ResponseWriter, r *http.Request) {
	var admissionReview admissionv1.AdmissionReview
	err := json.NewDecoder(r.Body).Decode(&admissionReview)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Delegate admission logic to the validation package
	response := validatePod(admissionReview.Request)

	// Populate AdmissionReview response
	admissionReview.Response = response
	admissionReview.Response.UID = admissionReview.Request.UID

	// Serialize the response and write it back
	resp, err := json.Marshal(admissionReview)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}
