resource "aws_eks_fargate_profile" "banking_kube_fargate_profile" {
  for_each               = toset(var.namespaces)
  fargate_profile_name   = each.key
  cluster_name           = var.cluster_name
  pod_execution_role_arn = var.fargate_pod_execution_role

  selector {
    namespace = each.key
  }
}