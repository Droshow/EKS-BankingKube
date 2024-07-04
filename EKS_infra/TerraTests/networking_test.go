package test

import (
    "testing"
    "github.com/gruntwork-io/terratest/modules/terraform"
    "github.com/stretchr/testify/assert"
)

func TestMainTerraform(t *testing.T) {
    terraformOptions := &terraform.Options{
        TerraformDir: "/Users/martin.drotar/Student/open_banking/EKS-BankingKube",
        Vars: map[string]interface{}{
            "cluster_name":           "example-cluster",
            "domain_name":            "example.com",
            "create_acm_certificate": false,
        },
    }

    defer terraform.Destroy(t, terraformOptions)

    // Initialize and apply the Terraform code, checking for errors.
    initAndApplyOutput, err := terraform.InitAndApplyE(t, terraformOptions)
    assert.NoError(t, err, "Terraform init and apply should not error out")
    assert.Contains(t, initAndApplyOutput, "Apply complete", "Terraform apply should be successful")

    // Validate the behavior of the modules.
    vpcID, err := terraform.OutputE(t, terraformOptions, "vpc_id")
    assert.NoError(t, err, "Fetching VPC ID should not error out")
    assert.NotEmpty(t, vpcID, "VPC ID should not be empty")
}