variable "certificate_names" {
  description = "Names for the ACM certificates to create"
  type        = map(string)
  default = {
    #put how many certs do you want to create but acm not usable for client vpn
    "acm_cert" = "devsbridge.com"
    # "client_cert" = "vpn-client.devsbridge.com"
  }
}
variable "cluster_name" {
  description = "The name of the EKS cluster"
  type        = string
}

variable "domain_name" {
  description = "The domain name to use for the ACM certificate"
  type        = string
}

variable "fetch_existing_certificates" {
  description = "Boolean to decide whether to fetch existing ACM certificates"
  type        = bool
  default     = true
}


variable "policies" {
  description = "List of policy ARNs to attach to the EKS role"
  type        = list(string)
  default = [
    "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy",
    "arn:aws:iam::aws:policy/AmazonEKSVPCResourceController"
  ]
}
variable "route_53cert_validation" {
  description = "The Route 53 DNS validation records for the ACM certificate"
  type = list(object({
    fqdn         = string
    record_name  = string
    record_type  = string
    record_value = string
  }))
}
variable "tags" {
  description = "A map of tags to add to the resources, with specific tags for each certificate"
  type        = map(map(string))
  default = {
    "server_cert" = {
      Environment = "Banking-Kube"
      Name        = "vpn-server-cert"
    },
    "client_cert" = {
      Environment = "Banking-Kube"
      Name        = "vpn-client-cert"
    }
  }
}


variable "vpc_id" {
  description = "The ID of the VPC where the EKS cluster and its resources will be created"
  type        = string
}

variable "aws_account_id" {
  description = "The AWS account ID"
  type        = string
}