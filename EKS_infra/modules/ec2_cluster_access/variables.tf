variable "instance_type" {
  description = "The instance type to use for the EC2 instance"
  type        = string
}

variable "subnet_id" {
  description = "The subnet ID where the EC2 instance will be deployed"
  type        = string
}

variable "security_group_id" {
  description = "The security group ID to associate with the EC2 instance"
  type        = string
}

variable "tags" {
  description = "A map of tags to assign to the resource"
  type        = map(string)
}

variable "aws_account_id" {
  description = "aws account"
  type        = string
}
variable "aws_region" {
  description = "AWS region"
  type        = string
}