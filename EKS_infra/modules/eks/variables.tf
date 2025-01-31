variable "cluster_name" {
  description = "The name of the EKS cluster"
  type        = string
}

variable "subnet_ids" {
  description = "A list of subnet IDs to deploy the EKS cluster in"
  type        = list(string)
}

variable "security_group_ids" {
  description = "A list of security group IDs to attach to the EKS cluster"
  type        = list(string)
}

variable "endpoint_private_access" {
  description = "Indicates whether or not to have private access enabled for the EKS cluster API server"
  type        = bool
  default     = true
}

variable "endpoint_public_access" {
  description = "Indicates whether or not to have public access enabled for the EKS cluster API server"
  type        = bool
  default     = false
}

variable "enabled_cluster_log_types" {
  description = "Types of EKS Cluster logging to enable"
  type        = list(string)
  default     = ["api", "audit"]
}

variable "cluster_role_iam_role_arn" {
  description = "The ARN of the IAM role to use for the EKS cluster"
  type        = string
}

variable "aws_account_id" {
  description = "The AWS account ID"
  type        = string
}