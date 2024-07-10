output "db_instance_address" {
  description = "The address of the RDS instance"
  value       = aws_db_instance.kube_db.address
}

output "db_instance_arn" {
  description = "The ARN of the RDS instance"
  value       = aws_db_instance.kube_db.arn
}

output "db_instance_name" {
  description = "The name of the RDS instance"
  value       = aws_db_instance.kube_db.identifier
}

output "db_instance_endpoint" {
  description = "The connection endpoint"
  value       = aws_db_instance.kube_db.endpoint
}

output "db_instance_hosted_zone_id" {
  description = "The canonical hosted zone ID of the DB instance (to be used in a Route 53 Alias record)"
  value       = aws_db_instance.kube_db.hosted_zone_id
}