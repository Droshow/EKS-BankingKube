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
#uncomment for self-hosted runners deployments

# provider "kubernetes" {
#   host                   = module.eks.cluster_endpoint
#   cluster_ca_certificate = module.eks.cluster_certificate_authority
#   exec {
#     api_version = "client.authentication.k8s.io/v1beta1"
#     command     = "aws"
#     args = [
#       "eks", "get-token",
#       "--cluster-name", module.eks.cluster_name,
#       "--region", var.aws_region
#     ]
#   }
# }