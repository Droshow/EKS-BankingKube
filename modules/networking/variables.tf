variable "subnets" {
  description = "Subnet configuration"
  default = {
    public1  = { cidr = "10.0.1.0/24", public = true }
    public2  = { cidr = "10.0.2.0/24", public = true }
    private1 = { cidr = "10.0.3.0/24", public = false }
    private2 = { cidr = "10.0.4.0/24", public = false }
  }
}

variable "vpc_cidr_block" {
  description = "CIDR block for the VPC"
  default     = "10.0.0.0/16"
}

variable "vpc_name" {
  description = "Name of the VPC"
  default     = "mother"
}

variable "security_groups" {
  description = "Security groups for the ALB"
}
variable "acm_domain_validation_options" {
  description = "Domain validation options of the ACM certificate"
  type        = any
}
variable "route53_zone_id" {
  description = "The ID of the Route53 hosted zone"
  type        = string
}
variable "vpc_id" {
  description = "The ID of the VPC"
  type        = string
}
variable "acm_certificate_arn" {
  description = "The ACM certificate"
  type        = any
}