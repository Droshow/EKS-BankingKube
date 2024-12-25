# provision ecr repository
resource "aws_ecr_repository" "ecr" {
  name = var.ecr_repo_name
  image_tag_mutability = "MUTABLE"
  tags = {
    Name = var.ecr_repo_name
  }
}