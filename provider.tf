provider "aws" {
  region = "eu-west-1"
  default_tags {
    tags = {
      project_name = "EKS-BankingKube"
      owner        = "DevsBridge"
      managed_by   = "Terraform"
    }
  }
}