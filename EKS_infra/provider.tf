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
# provider "kubernetes" {
#   host                   = module.eks.cluster_endpoint
#   cluster_ca_certificate = module.eks.cluster_certificate_authority_decoded
#   token                  = data.aws_eks_cluster_auth.auth.token
# }