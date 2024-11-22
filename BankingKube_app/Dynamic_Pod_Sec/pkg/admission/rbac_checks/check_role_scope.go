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

// RoleScopePolicy defines the structure for role scope policies
type RoleScopePolicy struct {
	RestrictedNamespaces []string `yaml:"restrictedNamespaces"`
}

// CheckRoleScope validates if a Role or ClusterRole complies with role scope policies
func CheckRoleScope(request *admissionv1.AdmissionRequest) bool {
	var role interface{}
	if request.Kind.Kind == "ClusterRole" {
		role = &rbacv1.ClusterRole{}
	} else if request.Kind.Kind == "Role" {
		role = &rbacv1.Role{}
	} else {
		log.Println("Unsupported kind:", request.Kind.Kind)
		return false
	}

	err := json.Unmarshal(request.Object.Raw, role)
	if err != nil {
		log.Println("Failed to parse RBAC role object:", err)
		return false
	}

	roleScopePolicy, err := getRoleScopePolicy()
	if err != nil {
		log.Println("Failed to load role scope policy:", err)
		return false
	}

	switch r := role.(type) {
	case *rbacv1.ClusterRole:
		if !validateRoleScope(r.Rules, roleScopePolicy) {
			log.Printf("ClusterRole %s has restricted scope\n", r.Name)
			return false
		}
	case *rbacv1.Role:
		if !validateRoleScope(r.Rules, roleScopePolicy) {
			log.Printf("Role %s has restricted scope\n", r.Name)
			return false
		}
	default:
		log.Println("Unsupported role type")
		return false
	}

	return true
}

// getRoleScopePolicy loads the role scope policy from the configuration file
func getRoleScopePolicy() (*RoleScopePolicy, error) {
	configPath := os.Getenv("SECURITY_POLICIES_PATH")
	if configPath == "" {
		configPath = "configs/security-policies.yaml"
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	var policies struct {
		RoleScope RoleScopePolicy `yaml:"roleScope"`
	}
	err = yaml.Unmarshal(data, &policies)
	return &policies.RoleScope, err
}

// validateRoleScope checks if the rules comply with the role scope policy
func validateRoleScope(rules []rbacv1.PolicyRule, policy *RoleScopePolicy) bool {
	for _, rule := range rules {
		for _, namespace := range rule.ResourceNames {
			if utils.Contains(policy.RestrictedNamespaces, namespace) {
				return false
			}
		}
	}
	return true
}
