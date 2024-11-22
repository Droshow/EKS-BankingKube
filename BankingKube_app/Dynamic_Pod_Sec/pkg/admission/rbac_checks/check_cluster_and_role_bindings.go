package rbac_checks

import (
	"encoding/json"
	"log"
	"os"

	"github.com/Droshow/EKS-BankingKube/BankingKube_app/Dynamic_Pod_Sec/pkg/utils"
	"gopkg.in/yaml.v2"
	admissionv1 "k8s.io/api/admission/v1"
	rbacv1 "k8s.io/api/rbac/v1"
)

// CheckRoleBindingsPolicy defines the structure for RBAC policies
type CheckRoleBindingsPolicy struct {
	RestrictedClusterRoles []string `yaml:"restrictedClusterRoles"`
	RestrictedRoleBindings []string `yaml:"restrictedRoleBindings"`
}

// CheckRBACBinding validates if a ClusterRoleBinding or RoleBinding complies with RBAC policies
func CheckRBACBinding(request *admissionv1.AdmissionRequest) bool {
	var roleBinding interface{}
	if request.Kind.Kind == "ClusterRoleBinding" {
		roleBinding = &rbacv1.ClusterRoleBinding{}
	} else if request.Kind.Kind == "RoleBinding" {
		roleBinding = &rbacv1.RoleBinding{}
	} else {
		log.Println("Unsupported kind:", request.Kind.Kind)
		return false
	}

	err := json.Unmarshal(request.Object.Raw, roleBinding)
	if err != nil {
		log.Println("Failed to parse RBAC binding object:", err)
		return false
	}

	rbacPolicy, err := getRBACPolicy()
	if err != nil {
		log.Println("Failed to load RBAC policy:", err)
		return false
	}

	switch binding := roleBinding.(type) {
	case *rbacv1.ClusterRoleBinding:
		if utils.Contains(rbacPolicy.RestrictedClusterRoles, binding.RoleRef.Name) {
			log.Printf("ClusterRoleBinding %s uses restricted ClusterRole %s\n", binding.Name, binding.RoleRef.Name)
			return false
		}
	case *rbacv1.RoleBinding:
		if utils.Contains(rbacPolicy.RestrictedRoleBindings, binding.RoleRef.Name) {
			log.Printf("RoleBinding %s uses restricted Role %s\n", binding.Name, binding.RoleRef.Name)
			return false
		}
	default:
		log.Println("Unsupported binding type")
		return false
	}

	return true
}

// getRBACPolicy loads the RBAC policy from the configuration file
func getRBACPolicy() (*CheckRoleBindingsPolicy, error) {
	configPath := os.Getenv("SECURITY_POLICIES_PATH")
	if configPath == "" {
		configPath = "configs/security-policies.yaml"
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	var policies struct {
		CheckRoleBindings CheckRoleBindingsPolicy `yaml:"checkRoleBindings"`
	}
	err = yaml.Unmarshal(data, &policies)
	return &policies.CheckRoleBindings, err
}
