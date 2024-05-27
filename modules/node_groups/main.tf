resource "aws_eks_fargate_profile" "default" {
  for_each = toset(var.namespaces)
  fargate_profile_name  = each.key
  cluster_name          = var.cluster_name
  pod_execution_role_arn = aws_iam_role.fargate_pod_execution_role.arn

  selector {
    namespace = each.key
  }
}