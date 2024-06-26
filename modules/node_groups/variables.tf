variable "cluster_name" {
  description = "The name of the EKS cluster"
  type        = string
}
# has to be automated with a controller that everytime a namespace is created an EKS Fargate profile is created as well.
# Another option is to Integrate this within CI/CD pipeline

variable "namespaces" {
  description = "List of namespaces for Fargate profiles"
  type        = list(string)
  default     = ["namespace1", "namespace2", "namespace3"]
}
variable "fargate_pod_execution_role" {
  description = "IAM role for Fargate pod execution"
}