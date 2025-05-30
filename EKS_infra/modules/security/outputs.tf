output "alb_sg_id" {
  description = "The ID of the security group for the ALB"
  value       = aws_security_group.alb_sg.id
}

output "acm_certificate_arn" {
  description = "The ARN of the ACM certificate"
  value       = aws_acm_certificate.cert["acm_cert"].arn
}

output "acm_certificate_arn_existing" {
  value = try(data.aws_acm_certificate.existing_cert["acm_cert"].arn, "")
}

output "domain_validation_options" {
  description = "Domain validation options of the ACM certificates"
  value = {
    "acm_cert" = [for option in aws_acm_certificate.cert["acm_cert"].domain_validation_options : {
      domain_name           = option.domain_name
      resource_record_name  = option.resource_record_name
      resource_record_value = option.resource_record_value
      resource_record_type  = option.resource_record_type
    }]
  }
  # value = {
  #   "server_cert" = [for option in aws_acm_certificate.cert["server_cert"].domain_validation_options : {
  #     domain_name           = option.domain_name
  #     resource_record_name  = option.resource_record_name
  #     resource_record_value = option.resource_record_value
  #     resource_record_type  = option.resource_record_type
  #   }]
  #   "client_cert" = [for option in aws_acm_certificate.cert["client_cert"].domain_validation_options : {
  #     domain_name           = option.domain_name
  #     resource_record_name  = option.resource_record_name
  #     resource_record_value = option.resource_record_value
  #     resource_record_type  = option.resource_record_type
  #   }]
  # }
}

output "eks_cluster_sg_id" {
  description = "The ID of the security group for the EKS cluster"
  value       = aws_security_group.eks_cluster_sg.id
}
output "eks_cluster_role_arn" {
  value = aws_iam_role.eks_cluster_role.arn
}

output "efs_security_group_id" {
  description = "The ID of the security group for EFS"
  value       = aws_security_group.efs_sg.id
}

output "alb_sg_arn" {
  description = "The ARN of the security group"
  value       = aws_security_group.alb_sg.arn
}
output "fargate_pod_execution_role_arn" {
  description = "The ARN of the IAM role for Fargate pod execution"
  value       = aws_iam_role.fargate_pod_execution_role.arn
}
output "worker_nodes_sg_id" {
  description = "The ID of the security group for the EKS worker nodes"
  value       = aws_security_group.worker_nodes_sg.id
}
output "ec2_access_aws_security_group" {
  description = "The ID of the security group for the EC2 instances"
  value       = aws_security_group.ec2_cluster_access_sg.id
}

output "vpc_endpoint_sg" {
  description = "vpc_endpoint_sg "
  value       = aws_security_group.vpc_endpoint_sg.id
}