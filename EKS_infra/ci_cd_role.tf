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
}

# Custom policy for additional permissions
resource "aws_iam_policy" "ci_cd_custom_policy" {
  name        = "ci-cd-custom-policy"
  description = "Custom policy for additional permissions required by CI/CD role"
  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Action = [
          "ec2:DescribeImages",
          "acm:ListCertificates"
        ],
        Resource = "*"
      }
    ]
  })
}

# Attach Policies to CI/CD Role
resource "aws_iam_role_policy_attachment" "ci_cd_role_policies" {
  for_each = toset([
    "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy",
    "arn:aws:iam::aws:policy/AmazonEKSServicePolicy",
    "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryFullAccess",
    "arn:aws:iam::aws:policy/AmazonS3FullAccess",
    "arn:aws:iam::aws:policy/IAMFullAccess",
    "arn:aws:iam::aws:policy/AmazonVPCFullAccess",
    "arn:aws:iam::aws:policy/AmazonEC2FullAccess",
    "arn:aws:iam::aws:policy/AWSCertificateManagerFullAccess",
    aws_iam_policy.ci_cd_custom_policy.arn # Attach custom policy
  ])
  role       = aws_iam_role.ci_cd_role.name
  policy_arn = each.value
}