variable "domain_name" {
  description = "The domain name to use for the ACM certificate"
  type        = string
  default     = "devsbridge.com"
}

variable "cluster_name" {
  description = "The name of the EKS cluster"
  type        = string
  default     = "Banking-Kube-Sloth"
}

variable "create_acm_certificate" {
  description = "Whether to create the ACM certificate"
  type        = bool
  default     = true
}
variable "aws_account_id" {
  description = "The AWS account ID"
  type        = string
  default     = "961477247679"
}

variable "github_repository" {
  description = "The GitHub repository in the format owner/repo"
  type        = string
}

variable "github_ref_name" {
  description = "The GitHub branch or tag name"
  type        = string
}