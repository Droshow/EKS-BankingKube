output "eks_cluster_sg_id" {
  description = "The ID of the security group for the EKS cluster"
  value       = aws_security_group.eks_cluster_sg.id
}

output "worker_nodes_sg_id" {
  description = "The ID of the security group for the EKS worker nodes"
  value       = aws_security_group.worker_nodes_sg.id
}

output "alb_sg_id" {
  description = "The ID of the security group for the ALB"
  value       = aws_security_group.alb_sg.id
}