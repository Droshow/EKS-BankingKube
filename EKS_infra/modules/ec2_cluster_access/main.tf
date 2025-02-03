locals {
  eks_policies = [
    "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore",
    "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy",
    "arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy",
    "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly",
    "arn:aws:iam::aws:policy/AmazonEC2FullAccess",
    "arn:aws:iam::aws:policy/AmazonVPCFullAccess",
    "arn:aws:iam::aws:policy/AmazonS3FullAccess"
  ]
}

data "aws_secretsmanager_secret" "github_runner" {
  name = "github_runner"
}
data "aws_secretsmanager_secret_version" "github_runner" {
  secret_id = data.aws_secretsmanager_secret.github_runner.id
}
data "aws_ami" "latest_amazon_linux" {
  most_recent = true
  owners      = ["amazon"] # Amazon's official AMIs

  filter {
    name   = "name"
    values = ["amzn2-ami-hvm-*-x86_64-gp2"] # Pattern for Amazon Linux 2 AMIs
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  filter {
    name   = "root-device-type"
    values = ["ebs"]
  }
}

resource "aws_instance" "ec2_cluster_access" {
  ami                         = data.aws_ami.latest_amazon_linux.id
  instance_type               = var.instance_type
  subnet_id                   = var.subnet_id
  associate_public_ip_address = false # this true only for public subnet
  vpc_security_group_ids      = [var.security_group_id]

  tags = var.tags

  iam_instance_profile = aws_iam_instance_profile.ec2_eks_profile.name

  user_data = <<-EOF
              #!/bin/bash
              yum install -y amazon-ssm-agent
              systemctl enable amazon-ssm-agent
              systemctl start amazon-ssm-agent

              # Install kubectl
              curl -o kubectl https://s3.us-west-2.amazonaws.com/amazon-eks/1.30.2/2024-07-12/bin/linux/amd64/kubectl
              chmod +x ./kubectl
              mv ./kubectl /usr/local/bin

              # Install aws-iam-authenticator
              curl -o aws-iam-authenticator https://amazon-eks.s3.us-west-2.amazonaws.com/1.21.2/2021-07-05/bin/linux/amd64/aws-iam-authenticator
              chmod +x ./aws-iam-authenticator
              mv ./aws-iam-authenticator /usr/local/bin

              # Remove old AWS CLI version 1 if present
              sudo rm -rf /usr/local/aws
              sudo rm /usr/bin/aws

               Install AWS CLI version 2.x
              sudo curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
              sudo unzip awscliv2.zip
              sudo ./aws/install

              # Create symlink for AWS CLI in /usr/local/bin
              sudo ln -s /usr/local/aws-cli/v2/current/bin/aws /usr/local/bin/aws

              # Verify installation
              aws --version
              kubectl version --client
              aws-iam-authenticator version

               # Fetch GitHub runner token from Secrets Manager just in case Terraform fails
              GITHUB_RUNNER_TOKEN=$(aws secretsmanager get-secret-value --secret-id github_runner --query SecretString --output text)
              echo "GitHub Runner Token: $GITHUB_RUNNER_TOKEN"

              # Install GitHub Actions Runner
              mkdir /home/ec2-user/actions-runner && cd /home/ec2-user/actions-runner
              curl -o actions-runner-linux-x64-2.321.0.tar.gz -L https://github.com/actions/runner/releases/download/v2.321.0/actions-runner-linux-x64-2.321.0.tar.gz
              echo "ba46ba7ce3a4d7236b16fbe44419fb453bc08f866b24f04d549ec89f1722a29e  actions-runner-linux-x64-2.321.0.tar.gz" | shasum -a 256 -c
              tar xzf ./actions-runner-linux-x64-2.321.0.tar.gz

              # Configure the GitHub Actions Runner use both commands with Terraform OR AWS Fetch to be sure 
              ./config.sh --url https://github.com/Droshow/EKS-BankingKube --token ${data.aws_secretsmanager_secret_version.github_runner.secret_string} --unattended --replace || ./config.sh --url https://github.com/Droshow/EKS-BankingKube --token $GITHUB_RUNNER_TOKEN --unattended --replace
              # Start the GitHub Actions Runner as a service
              ./svc.sh install
              ./svc.sh start

              
              EOF
}


resource "aws_iam_role" "ec2_eks_role" {
  name = "ec2-eks-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "ec2.amazonaws.com"
        }
      }
    ]
  })
}

# Attach all required policies in a loop
resource "aws_iam_role_policy_attachment" "eks_policies" {
  for_each = toset(local.eks_policies)

  role       = aws_iam_role.ec2_eks_role.name
  policy_arn = each.value
}

# Instance profile for the EC2 role
resource "aws_iam_instance_profile" "ec2_eks_profile" {
  name = "ec2-eks-instance-profile"
  role = aws_iam_role.ec2_eks_role.name
}