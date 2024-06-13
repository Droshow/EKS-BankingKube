package test

import (
    "testing"
    "github.com/gruntwork-io/terratest/modules/terraform"
)

func TestNetworking(t *testing.T) {
    terraformOptions := &terraform.Options{
        // The path to where your Terraform code is located
        TerraformDir: "../modules/networking",
    }

    // This will run `terraform init` and `terraform apply` and fail the test if there are any errors
    terraform.InitAndApply(t, terraformOptions)

    // At the end of the test, run `terraform destroy` to clean up any resources that were created
    defer terraform.Destroy(t, terraformOptions)
}