variable "allocated_storage" {
  description = "The allocated storage in gibibytes."
  type        = number
  default     = 20
}

variable "storage_type" {
  description = "One of standard (magnetic), gp2 (general purpose SSD), or io1 (provisioned IOPS SSD)."
  type        = string
  default     = "gp2"
}

variable "engine" {
  description = "The name of the database engine to be used for this DB instance."
  type        = string
  default     = "mysql"
}

variable "engine_version" {
  description = "The version number of the database engine to use."
  type        = string
  default     = "5.7"
}

variable "instance_class" {
  description = "The compute and memory capacity of the DB instance."
  type        = string
  default     = "db.t2.micro"
}

variable "db_name" {
  description = "The name of the database to create when the DB instance is created."
  type        = string
}

variable "username" {
  description = "The name of master user for the client DB instance."
  type        = string
}

variable "password" {
  description = "The password for the master database user."
  type        = string
}

variable "parameter_group_name" {
  description = "The name of the DB parameter group to associate with this DB instance."
  type        = string
  default     = "default.mysql5.7"
}

variable "skip_final_snapshot" {
  description = "Determines whether a final DB snapshot is created before the DB instance is deleted."
  type        = bool
  default     = true
}

variable "db_subnets" {
  description = "A list of VPC subnet IDs."
  type        = list(string)
}
