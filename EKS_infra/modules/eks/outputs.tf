output "cluster_endpoint" {
  description = "The endpoint for your EKS cluster"
  value       = aws_eks_cluster.banking_kube_cluster.endpoint
}

# output "kubeconfig_certificate_authority_data" {
#   description = "The certificate authority data for your EKS cluster"
#   value       = aws_eks_cluster.banking_kube_cluster.certificate_authority[0].data
# }

output "cluster_name" {
  description = "The name of the EKS cluster"
  value       = aws_eks_cluster.banking_kube_cluster.name
}

output "cluster_arn" {
  description = "The ARN of the EKS cluster"
  value       = aws_eks_cluster.banking_kube_cluster.arn
}
output "cluster_certificate_authority" {
  description = "The certificate authority for the EKS cluster"
  value       = base64decode(aws_eks_cluster.banking_kube_cluster.certificate_authority[0].data)
}
