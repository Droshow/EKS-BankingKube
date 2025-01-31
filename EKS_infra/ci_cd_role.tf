################################
# Create the CI/CD Role
################################
resource "aws_iam_role" "ci_cd_role" {
  name = "ci-cd-role"
  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Principal = {
          Federated = "arn:aws:iam::${var.aws_account_id}:oidc-provider/token.actions.githubusercontent.com"
        },
        Action = "sts:AssumeRoleWithWebIdentity",
        Condition = {
          StringLike = {
            "token.actions.githubusercontent.com:sub" = "repo:Droshow/EKS-BankingKube:*",
            "token.actions.githubusercontent.com:aud" = "sts.amazonaws.com"
          }
        }
      },
    ]
  })
  lifecycle {
    prevent_destroy = true
  }
}

################################
# Create the Custom Policy
################################
resource "aws_iam_policy" "ci_cd_custom_policy" {
  name        = "ci-cd-custom-policy"
  description = "Custom policy for additional permissions required by CI/CD role"
  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Action = [
          "ec2:*",
          "acm:*",
          "secretsmanager:*",
          "rds:*",
          "eks:*",
          "elasticloadbalancing:*",
          "route53:*",
          "elasticfilesystem:*"
        ],
        Resource = "*"
      }
    ]
  })
}

################################
# Attach Known AWS-Managed Policies
################################
# 1) Put your known ARNs in a local variable
locals {
  ci_cd_known_arns = [
    "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy",
    "arn:aws:iam::aws:policy/AmazonEKSServicePolicy",
    "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryFullAccess",
    "arn:aws:iam::aws:policy/AmazonS3FullAccess",
    "arn:aws:iam::aws:policy/IAMFullAccess",
    "arn:aws:iam::aws:policy/AmazonVPCFullAccess",
    "arn:aws:iam::aws:policy/AmazonEC2FullAccess",
    "arn:aws:iam::aws:policy/AWSCertificateManagerFullAccess"
  ]
}

resource "aws_iam_role_policy_attachment" "known_policies_attach" {
  for_each   = toset(local.ci_cd_known_arns)
  role       = aws_iam_role.ci_cd_role.name
  policy_arn = each.value
}

################################
# Attach the Custom Policy ARN
################################
# 2) Separate resource for the custom policy since it's unknown until apply
resource "aws_iam_role_policy_attachment" "custom_policy_attach" {
  role       = aws_iam_role.ci_cd_role.name
  policy_arn = aws_iam_policy.ci_cd_custom_policy.arn
}
