output "instance_id" {
  description = "The ID of the EC2 instance"
  value       = aws_instance.ec2_cluster_access.id
}
output "instance_private_ip" {
  description = "The private IP address of the EC2 instance"
  value       = aws_instance.ec2_cluster_access.private_ip
}

output "instance_arn" {
  description = "The ARN of the EC2 instance"
  value       = aws_instance.ec2_cluster_access.arn
}