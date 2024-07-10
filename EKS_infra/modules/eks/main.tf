resource "aws_eks_cluster" "banking_kube_cluster" {
  name     = var.cluster_name
  role_arn = var.cluster_role_iam_role_arn

  vpc_config {
    subnet_ids              = var.subnet_ids
    endpoint_private_access = var.endpoint_private_access
    endpoint_public_access  = var.endpoint_public_access
    security_group_ids      = var.security_group_ids
  }
  access_config {
    authentication_mode                         = "API_AND_CONFIG_MAP"
    bootstrap_cluster_creator_admin_permissions = true # false in later chapters
  }

  enabled_cluster_log_types = var.enabled_cluster_log_types
}
