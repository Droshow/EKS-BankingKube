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

#########
#config map
# data "kubernetes_config_map" "aws_auth" {
#   metadata {
#     name      = "aws-auth"
#     namespace = "kube-system"
#   }
# }

# resource "kubernetes_config_map" "aws_auth" {
#   metadata {
#     name      = "aws-auth"
#     namespace = "kube-system"
#   }

#   data = {
#     mapRoles = yamlencode([
#       for role in concat(
#         [
#           {
#             rolearn  = "arn:aws:iam::${var.aws_account_id}:role/ci-cd-role"
#             username = "admin"
#             groups   = ["system:masters"]
#           },
#           {
#             rolearn  = "arn:aws:iam::${var.aws_account_id}:role/fargate-pod-execution-role"
#             username = "fargate-pod-execution-role"
#             groups   = ["system:node:{{SessionName}}", "system:nodes", "system:node-proxier"]
#           },
#           {
#             rolearn  = "arn:aws:iam::${var.aws_account_id}:role/ec2-eks-role"
#             username = "ec2-eks-role"
#             groups   = ["system:masters"]
#           }
#         ],
#         [
#           for role in yamldecode(data.kubernetes_config_map.aws_auth.data["mapRoles"]) : role
#         ]
#       ) : role
#     ])
#   }
# }