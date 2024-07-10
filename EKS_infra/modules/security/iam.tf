data "aws_iam_policy_document" "assume_role" {
  statement {
    effect = "Allow"
    principals {
      type        = "Service"
      identifiers = ["eks.amazonaws.com"]
    }
    actions = ["sts:AssumeRole"]
  }
}

resource "aws_iam_role" "alb_ingress_controller" {
  name = "alb-ingress-controller"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "eks.amazonaws.com"
        }
      },
    ]
  })
}

resource "aws_iam_role_policy_attachment" "alb_ingress_controller" {
  role       = aws_iam_role.alb_ingress_controller.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"
}

resource "aws_iam_role" "eks_cluster_role" {
  name               = "${var.cluster_name}-eks-role"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
}

resource "aws_iam_role_policy_attachment" "eks_policies" {
  for_each   = toset(var.policies)
  policy_arn = each.key
  role       = aws_iam_role.eks_cluster_role.name
}
resource "aws_iam_role" "fargate_pod_execution_role" {
  name = "fargate-pod-execution-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "eks-fargate-pods.amazonaws.com"
        }
      },
    ]
  })
}


