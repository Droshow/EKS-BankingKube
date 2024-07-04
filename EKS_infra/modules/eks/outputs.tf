output "cluster_endpoint" {
  description = "The endpoint for your EKS cluster"
  value       = aws_eks_cluster.example.endpoint
}

output "kubeconfig_certificate_authority_data" {
  description = "The certificate authority data for your EKS cluster"
  value       = aws_eks_cluster.example.certificate_authority[0].data
}
output "cluster_name" {
  description = "The name of the EKS cluster"
  value       = var.cluster_name
}