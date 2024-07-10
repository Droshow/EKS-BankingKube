output "alb_sg_id" {
  description = "The ID of the security group for the ALB"
  value       = aws_security_group.alb_sg.id
}
output "certificate_arn" {
  description = "The ARN of the ACM certificate"
  value       = length(aws_acm_certificate.cert) > 0 ? aws_acm_certificate.cert[0].arn : null
}

output "domain_validation_options" {
  description = "Domain validation options of the ACM certificate"
  value = length(aws_acm_certificate.cert) > 0 ? [for option in aws_acm_certificate.cert[0].domain_validation_options : {
    domain_name           = option.domain_name
    resource_record_name  = option.resource_record_name
    resource_record_value = option.resource_record_value
    resource_record_type  = option.resource_record_type
  }] : []
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
output "eks_alb_sg_id" {
  description = "The ID of the security group"
  value       = aws_security_group.eks_alb_sg.id
}
output "eks_alb_sg_arn" {
  description = "The ARN of the security group"
  value       = aws_security_group.eks_alb_sg.arn
}
output "fargate_pod_execution_role_arn" {
  description = "The ARN of the IAM role for Fargate pod execution"
  value       = aws_iam_role.fargate_pod_execution_role.arn
}
output "validation_arns" {
  description = "The ARNs of the ACM certificate validations"
  value       = aws_acm_certificate_validation.cert.*.validation_record_fqdns
}

output "worker_nodes_sg_id" {
  description = "The ID of the security group for the EKS worker nodes"
  value       = aws_security_group.worker_nodes_sg.id
}
