output "server_certificate_arn" {
  description = "The ARN of the server ACM certificate"
  value       = var.server_certificate_arn
}

output "client_root_certificate_arn" {
  description = "The ARN of the client root ACM certificate"
  value       = var.client_root_certificate_arn
}

output "subnet_id" {
  description = "The ID of the subnet"
  value       = var.subnet_id
}