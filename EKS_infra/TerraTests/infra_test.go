package test

import (
    "testing"
    "github.com/gruntwork-io/terratest/modules/terraform"
    "github.com/stretchr/testify/assert"
)

func TestInfrastructureModules(t *testing.T) {
    terraformOptions := &terraform.Options{
        TerraformDir: "/Users/martin.drotar/Student/open_banking/EKS-BankingKube",
    }

    defer terraform.Destroy(t, terraformOptions)

    // Initialize and apply the Terraform code, checking for errors.
    initAndApplyOutput, err := terraform.InitAndApplyE(t, terraformOptions)
    assert.NoError(t, err, "Terraform init and apply should not error out")
    assert.Contains(t, initAndApplyOutput, "Apply complete", "Terraform apply should be successful")

    // Validate the behavior of the modules.
    // Networking Module
    vpcID, err := terraform.OutputE(t, terraformOptions, "vpc_id")
    assert.NoError(t, err, "Fetching VPC ID should not error out")
    assert.NotEmpty(t, vpcID, "VPC ID should not be empty")

    // Security Module
    albSgID, err := terraform.OutputE(t, terraformOptions, "alb_sg_id")
    assert.NoError(t, err, "Fetching ALB Security Group ID should not error out")
    assert.NotEmpty(t, albSgID, "ALB Security Group ID should not be empty")

    // EKS Module
    eksClusterID, err := terraform.OutputE(t, terraformOptions, "eks_cluster_id")
    assert.NoError(t, err, "Fetching EKS Cluster ID should not error out")
    assert.NotEmpty(t, eksClusterID, "EKS Cluster ID should not be empty")

    // DB Instance Module
    dbInstanceEndpoint, err := terraform.OutputE(t, terraformOptions, "db_instance_endpoint")
    assert.NoError(t, err, "Fetching DB Instance Endpoint should not error out")
    assert.NotEmpty(t, dbInstanceEndpoint, "DB Instance Endpoint should not be empty")

    // Storage Module
    efsID, err := terraform.OutputE(t, terraformOptions, "efs_id")
    assert.NoError(t, err, "Fetching EFS ID should not error out")
    assert.NotEmpty(t, efsID, "EFS ID should not be empty")
}