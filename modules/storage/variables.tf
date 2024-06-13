variable "subnet_ids" {
  description = "The IDs of the subnets where the EFS mount targets will be created"
  type        = list(string)
}
variable "vpc_id" {
  description = "The ID of the VPC where the EFS mount targets will be created"
  type        = string
}
variable "tags" {
  description = "A map of tags to add to all resources"
  type        = map(string)
}
variable "efs_sg_id" {
  description = "The security group for the EFS mount targets"
  type        = string
}