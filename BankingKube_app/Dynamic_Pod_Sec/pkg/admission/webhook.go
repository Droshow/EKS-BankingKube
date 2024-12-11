package admission

import (
	"encoding/json"
	"net/http"

	"github.com/Droshow/EKS-BankingKube/BankingKube_app/Dynamic_Pod_Sec/pkg/admission/api_restrictions"
	"github.com/Droshow/EKS-BankingKube/BankingKube_app/Dynamic_Pod_Sec/pkg/admission/context_capabilities"
	"github.com/Droshow/EKS-BankingKube/BankingKube_app/Dynamic_Pod_Sec/pkg/admission/image_security"
	"github.com/Droshow/EKS-BankingKube/BankingKube_app/Dynamic_Pod_Sec/pkg/admission/network_security"
	"github.com/Droshow/EKS-BankingKube/BankingKube_app/Dynamic_Pod_Sec/pkg/admission/rbac_checks"
	"github.com/Droshow/EKS-BankingKube/BankingKube_app/Dynamic_Pod_Sec/pkg/admission/resource_limits"
	"github.com/Droshow/EKS-BankingKube/BankingKube_app/Dynamic_Pod_Sec/pkg/admission/volume_security"

	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// HandleAdmissionRequest handles incoming admission requests based on the URL path
func HandleAdmissionRequest(w http.ResponseWriter, r *http.Request) {
	var admissionReview admissionv1.AdmissionReview
	err := json.NewDecoder(r.Body).Decode(&admissionReview)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Determine which validation to perform based on the request path
	var response *admissionv1.AdmissionResponse
	switch r.URL.Path {
	case "/validate/context":
		response = validateContext(admissionReview.Request)
	case "/validate/volumes":
		response = validateVolumes(admissionReview.Request)
	case "/validate/network":
		response = validateNetwork(admissionReview.Request)
	case "/validate/api":
		response = validateAPI(admissionReview.Request)
	case "/validate/image":
		response = validateImage(admissionReview.Request)
	case "/validate/rbac":
		response = validateRBAC(admissionReview.Request)
	case "/validate/resources":
		response = validateResources(admissionReview.Request)
	case "/mutate/pod":
		response = mutatePod(admissionReview.Request)
	default:
		http.Error(w, "Invalid validation path", http.StatusNotFound)
		return
	}

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

// validateContext handles context and capabilities checks
func validateContext(request *admissionv1.AdmissionRequest) *admissionv1.AdmissionResponse {
	allowed := true
	result := &metav1.Status{Message: "Pod context validation passed"}

	if !context_capabilities.CheckPodSecurityContext(request) {
		allowed = false
		result = &metav1.Status{Message: "Pod contains privileged containers, which is not allowed."}
	}
	if !context_capabilities.CheckCapabilities(request) {
		allowed = false
		result = &metav1.Status{Message: "Pod contains containers with disallowed capabilities."}
	}

	return &admissionv1.AdmissionResponse{Allowed: allowed, Result: result}
}

// validateVolumes handles volume security checks
func validateVolumes(request *admissionv1.AdmissionRequest) *admissionv1.AdmissionResponse {
	allowed := true
	result := &metav1.Status{Message: "Pod volume validation passed"}

	if !volume_security.CheckHostPath(request) {
		allowed = false
		result = &metav1.Status{Message: "Pod contains disallowed host paths."}
	}

	return &admissionv1.AdmissionResponse{Allowed: allowed, Result: result}
}

// validateNetwork handles network security checks
func validateNetwork(request *admissionv1.AdmissionRequest) *admissionv1.AdmissionResponse {
	allowed := true
	result := &metav1.Status{Message: "Pod network validation passed"}

	if !network_security.CheckPolicyConsistency() {
		allowed = false
		result = &metav1.Status{Message: "Security policies are inconsistent."}
	}
	if !network_security.CheckNetworkPolicy(request) {
		allowed = false
		result = &metav1.Status{Message: "Pod does not comply with network policies."}
	}
	if !network_security.CheckHostNetwork(request) {
		allowed = false
		result = &metav1.Status{Message: "Pod is using the host network, which is disallowed."}
	}
	if !network_security.CheckEgress(request) {
		allowed = false
		result = &metav1.Status{Message: "Pod has an egress route that violates policy."}
	}
	if !network_security.CheckIngress(request) {
		allowed = false
		result = &metav1.Status{Message: "Pod has an ingress route that violates policy."}
	}

	return &admissionv1.AdmissionResponse{Allowed: allowed, Result: result}
}

// validateAPI handles API access and service account checks
func validateAPI(request *admissionv1.AdmissionRequest) *admissionv1.AdmissionResponse {
	allowed := true
	result := &metav1.Status{Message: "Pod API access validation passed"}

	if !api_restrictions.CheckAPIAccess(request) {
		allowed = false
		result = &metav1.Status{Message: "Pod is attempting to access restricted API paths."}
	}
	if !api_restrictions.CheckServiceAccount(request) {
		allowed = false
		result = &metav1.Status{Message: "Pod is using a restricted service account."}
	}

	return &admissionv1.AdmissionResponse{Allowed: allowed, Result: result}
}

// validateImage handles image security checks
func validateImage(request *admissionv1.AdmissionRequest) *admissionv1.AdmissionResponse {
	allowed := true
	result := &metav1.Status{Message: "Pod image validation passed"}

	if !image_security.CheckImageRegistry(request) {
		allowed = false
		result = &metav1.Status{Message: "Pod is using an image from a disallowed registry."}
	}
	if !image_security.CheckImageSigning(request) {
		allowed = false
		result = &metav1.Status{Message: "Pod is using an unsigned image."}
	}
	if !image_security.CheckImageTags(request) {
		allowed = false
		result = &metav1.Status{Message: "Pod is using an image with a disallowed tag."}
	}

	return &admissionv1.AdmissionResponse{Allowed: allowed, Result: result}
}

// validateRBAC handles RBAC-related checks
func validateRBAC(request *admissionv1.AdmissionRequest) *admissionv1.AdmissionResponse {
	allowed := true
	result := &metav1.Status{Message: "RBAC validation passed"}

	if !rbac_checks.CheckRBACBinding(request) {
		allowed = false
		result = &metav1.Status{Message: "ClusterRoleBinding uses a restricted ClusterRole."}
	}
	if !rbac_checks.CheckPermissionLevels(request) {
		allowed = false
		result = &metav1.Status{Message: "Role or ClusterRole has restricted permissions."}
	}
	if !rbac_checks.CheckRoleScope(request) {
		allowed = false
		result = &metav1.Status{Message: "Role or ClusterRole has restricted scope."}
	}

	return &admissionv1.AdmissionResponse{Allowed: allowed, Result: result}
}

// validateResources handles resource limit and request checks
func validateResources(request *admissionv1.AdmissionRequest) *admissionv1.AdmissionResponse {
	allowed := true
	result := &metav1.Status{Message: "Pod resource validation passed"}

	if !resource_limits.CheckResourceLimits(request) {
		allowed = false
		result = &metav1.Status{Message: "Pod contains containers without resource limits."}
	}
	if !resource_limits.CheckResourceRequests(request) {
		allowed = false
		result = &metav1.Status{Message: "Pod contains containers without resource requests."}
	}

	return &admissionv1.AdmissionResponse{Allowed: allowed, Result: result}
}

// mutatePod applies baseline security configurations
func mutatePod(request *admissionv1.AdmissionRequest) *admissionv1.AdmissionResponse {
	pod := &corev1.Pod{}
	if err := json.Unmarshal(request.Object.Raw, pod); err != nil {
		return &admissionv1.AdmissionResponse{
			Allowed: false,
			Result:  &metav1.Status{Message: "Failed to parse Pod object"},
		}
	}

	applyBaselineSecurity(pod)

	patchBytes, err := json.Marshal(pod)
	if err != nil {
		return &admissionv1.AdmissionResponse{
			Allowed: false,
			Result:  &metav1.Status{Message: "Failed to generate patch"},
		}
	}

	return &admissionv1.AdmissionResponse{
		Allowed: true,
		Patch:   patchBytes,
		PatchType: func() *admissionv1.PatchType {
			pt := admissionv1.PatchTypeJSONPatch
			return &pt
		}(),
	}
}

// applyBaselineSecurity applies essential security defaults
func applyBaselineSecurity(pod *corev1.Pod) {
	for i := range pod.Spec.Containers {
		if pod.Spec.Containers[i].SecurityContext == nil {
			pod.Spec.Containers[i].SecurityContext = &corev1.SecurityContext{}
		}

		// Apply essential security defaults
		pod.Spec.Containers[i].SecurityContext.RunAsNonRoot = boolPtr(true)
		pod.Spec.Containers[i].SecurityContext.ReadOnlyRootFilesystem = boolPtr(true)
		pod.Spec.Containers[i].SecurityContext.AllowPrivilegeEscalation = boolPtr(false)

		// Drop disallowed capabilities
		pod.Spec.Containers[i].SecurityContext.Capabilities = &corev1.Capabilities{
			Drop: []corev1.Capability{"CAP_SYS_ADMIN", "CAP_NET_ADMIN"},
		}
	}
}

// Helper function to create boolean pointers
func boolPtr(b bool) *bool {
	return &b
}
