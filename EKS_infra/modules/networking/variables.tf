variable "subnets" {
  description = "Subnet configuration"
  type = map(object({
    cidr   = string
    public = bool
    name   = string
    az     = string
  }))
  default = {
    eks_public_subnet-001  = { cidr = "10.0.1.0/24", public = true, name = "eks_public_subnet-001", az = "eu-central-1a" },
    eks_public_subnet-002  = { cidr = "10.0.2.0/24", public = true, name = "eks_public_subnet-002", az = "eu-central-1b" },
    eks_private_subnet-003 = { cidr = "10.0.3.0/24", public = false, name = "eks_private_subnet-003", az = "eu-central-1a" },
    eks_private_subnet-004 = { cidr = "10.0.4.0/24", public = false, name = "eks_private_subnet-004", az = "eu-central-1b" }
  }
}

variable "vpc_cidr_block" {
  description = "CIDR block for the VPC"
  default     = "10.0.0.0/16"
}

variable "vpc_name" {
  description = "Name of the VPC"
  default     = "Banking-Kube-main"
}

variable "alb_security_group" {
  description = "Security groups for the ALB"
}
variable "acm_domain_validation_options" {
  description = "Domain validation options of the ACM certificates"
  type = map(list(object({
    domain_name           = string
    resource_record_name  = string
    resource_record_value = string
    resource_record_type  = string
  })))
}

variable "acm_certificate_arn" {
  description = "The ARN of the server ACM certificate"
  type        = string
}