variable "domain_name" {
  description = "The domain name to use for the ACM certificate"
  type        = string
  default     = "bankingkube.com"
}

variable "cluster_name" {
  description = "The name of the EKS cluster"
  type        = string
  default     = "Banking-Kube-Sloth"
}