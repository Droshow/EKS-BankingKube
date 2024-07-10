resource "aws_db_instance" "kube_db" {
  allocated_storage    = var.allocated_storage
  storage_type         = var.storage_type
  engine               = var.engine
  engine_version       = var.engine_version
  instance_class       = var.instance_class
  identifier           = var.db_name
  username             = var.username
  password             = var.password
  parameter_group_name = var.parameter_group_name
  skip_final_snapshot  = var.skip_final_snapshot
  db_subnet_group_name = aws_db_subnet_group.kube_db_subnet_group.name

}

resource "aws_db_subnet_group" "kube_db_subnet_group" {
  name       = "kube-db-subnet-group"
  subnet_ids = var.db_subnets

  tags = {
    Name = "Kube-Banking-Subnet-Group"
  }
}