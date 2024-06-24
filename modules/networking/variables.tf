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
  default     = "Banking-Kube-main"
}

variable "alb_security_group" {
  description = "Security groups for the ALB"
}
variable "acm_domain_validation_options" {
  description = "Domain validation options of the ACM certificate"
  type        = any
}
variable "acm_certificate_arn" {
  description = "The ACM certificate"
  type        = any
}