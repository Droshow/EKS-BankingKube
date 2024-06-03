variable "cluster_name" {
  description = "The name of the EKS cluster"
  type        = string
}

variable "vpc_id" {
  description = "The ID of the VPC where the EKS cluster and its resources will be created"
  type        = string
}

variable "policies" {
  description = "List of policy ARNs to attach to the EKS role"
  type        = list(string)
  default = [
    "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy",
    "arn:aws:iam::aws:policy/AmazonEKSVPCResourceController"
  ]
}
variable "domain_name" {
  description = "The domain name to use for the ACM certificate"
  type        = string
}

variable "tags" {
  description = "A map of tags to add to the resources"
  type        = map(string)
  default = {
    Environment = "production"
  }
}
