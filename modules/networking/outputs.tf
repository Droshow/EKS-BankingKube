output "vpc_id" {
  description = "The ID of the VPC"
  value       = aws_vpc.eks_vpc.id
}

output "aws_route_53_cert_validation" {
  description = "The Route 53 record for ACM certificate validation"
  value       = [for record in aws_route53_record.cert_validation : record.fqdn]
}