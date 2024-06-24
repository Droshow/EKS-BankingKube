package test

import (
    "testing"
    "github.com/gruntwork-io/terratest/modules/terraform"
    "github.com/stretchr/testify/assert"
)

func TestMainTerraform(t *testing.T) {
    terraformOptions := &terraform.Options{
        // Set the path to the Terraform code that will be tested.
        TerraformDir: "/Users/martin.drotar/Student/open_banking/EKS-BankingKube",

        // Variables to pass to our Terraform code using -var options
        Vars: map[string]interface{}{
            "cluster_name": "example-cluster",
            "domain_name":  "example.com",
        },
    }

    // Clean up resources with "terraform destroy" at the end of the test.
    defer terraform.Destroy(t, terraformOptions)

    // Initialize and apply the Terraform code.
    terraform.InitAndApply(t, terraformOptions)

    // Add assertions here to validate the behavior of the modules.
    // For example, assert that the VPC ID is not empty.
    vpcID := terraform.Output(t, terraformOptions, "vpc_id")
    assert.NotEmpty(t, vpcID, "VPC ID should not be empty")
}