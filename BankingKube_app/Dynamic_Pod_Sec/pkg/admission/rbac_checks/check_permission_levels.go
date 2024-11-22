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

// PermissionLevelsPolicy defines the structure for permission level policies
type PermissionLevelsPolicy struct {
	RestrictedVerbs     []string `yaml:"restrictedVerbs"`
	RestrictedResources []string `yaml:"restrictedResources"`
}

// CheckPermissionLevels validates if a Role or ClusterRole complies with permission level policies
func CheckPermissionLevels(request *admissionv1.AdmissionRequest) bool {
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

	permissionPolicy, err := getPermissionLevelsPolicy()
	if err != nil {
		log.Println("Failed to load permission levels policy:", err)
		return false
	}

	switch r := role.(type) {
	case *rbacv1.ClusterRole:
		if !validateRules(r.Rules, permissionPolicy) {
			log.Printf("ClusterRole %s has restricted permissions\n", r.Name)
			return false
		}
	case *rbacv1.Role:
		if !validateRules(r.Rules, permissionPolicy) {
			log.Printf("Role %s has restricted permissions\n", r.Name)
			return false
		}
	default:
		log.Println("Unsupported role type")
		return false
	}

	return true
}

// getPermissionLevelsPolicy loads the permission levels policy from the configuration file
func getPermissionLevelsPolicy() (*PermissionLevelsPolicy, error) {
	configPath := os.Getenv("SECURITY_POLICIES_PATH")
	if configPath == "" {
		configPath = "configs/security-policies.yaml"
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	var policies struct {
		PermissionLevels PermissionLevelsPolicy `yaml:"permissionLevels"`
	}
	err = yaml.Unmarshal(data, &policies)
	return &policies.PermissionLevels, err
}

// validateRules checks if the rules comply with the permission levels policy
func validateRules(rules []rbacv1.PolicyRule, policy *PermissionLevelsPolicy) bool {
	for _, rule := range rules {
		for _, verb := range rule.Verbs {
			if utils.Contains(policy.RestrictedVerbs, verb) {
				return false
			}
		}
		for _, resource := range rule.Resources {
			if utils.Contains(policy.RestrictedResources, resource) {
				return false
			}
		}
	}
	return true
}
