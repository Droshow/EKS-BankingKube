output "vpc_id" {
  description = "The ID of the VPC"
  value       = aws_vpc.eks_vpc.id
}

output "aws_route_53_cert_validation" {
  description = "The Route 53 DNS validation records for the ACM certificate"
  value = [for dvo in var.acm_domain_validation_options : {
    fqdn         = "${dvo.resource_record_name}.${aws_route53_zone.banking-kube.name}"
    record_name  = dvo.resource_record_name
    record_type  = dvo.resource_record_type
    record_value = dvo.resource_record_value
  }]
}

output "private_subnets_ids" {
  description = "The IDs of the private subnets"
  value       = [for k in keys(locals.private_subnets) : aws_subnet.subnet[k].id]
}