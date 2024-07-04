provider "aws" {
  region = "eu-central-1"
  default_tags {
    tags = {
      project_name = "EKS-BankingKube"
      owner        = "DevsBridge"
      managed_by   = "Terraform"
    }
  }
}