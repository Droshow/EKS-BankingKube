locals {
  private_subnets = { for k, v in var.subnets : k => v if !v.public }
  public_subnets  = { for k, v in var.subnets : k => v if v.public }
  vpc_endpoints = {
    ssm          = "com.amazonaws.eu-central-1.ssm"
    ssm_messages = "com.amazonaws.eu-central-1.ssmmessages"
    ec2_messages = "com.amazonaws.eu-central-1.ec2messages"
  }

}
#Hello VPC
resource "aws_vpc" "eks_vpc" {
  cidr_block           = var.vpc_cidr_block
  enable_dns_support   = true
  enable_dns_hostnames = true
}

resource "aws_subnet" "subnet" {
  for_each = var.subnets

  vpc_id                  = aws_vpc.eks_vpc.id
  cidr_block              = each.value.cidr
  availability_zone       = each.value.az
  map_public_ip_on_launch = each.value.public
  tags = {
    "Name" = "${each.value.name} (${each.value.public ? "Public" : "Private"})",
    "AZ"   = each.value.az

  }
}

resource "aws_route_table" "public" {
  vpc_id = aws_vpc.eks_vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.gw.id
  }

  tags = {
    Name = "public"
  }
}

resource "aws_route_table" "private" {
  vpc_id = aws_vpc.eks_vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_nat_gateway.ipam[keys(local.public_subnets)[0]].id # Use the first NAT Gateway
  }

  tags = {
    Name = "private"
  }
}

resource "aws_internet_gateway" "gw" {
  vpc_id = aws_vpc.eks_vpc.id
}
resource "aws_route_table_association" "public" {
  for_each = { for name, subnet in var.subnets : name => subnet if subnet.public }

  subnet_id      = aws_subnet.subnet[each.key].id
  route_table_id = aws_route_table.public.id
}

resource "aws_route_table_association" "private" {
  for_each = { for name, subnet in var.subnets : name => subnet if !subnet.public }

  subnet_id      = aws_subnet.subnet[each.key].id
  route_table_id = aws_route_table.private.id
}

resource "aws_eip" "ipam" {
  for_each = { for name, subnet in var.subnets : name => subnet if subnet.public }
}

resource "aws_nat_gateway" "ipam" {
  for_each = { for name, subnet in var.subnets : name => subnet if subnet.public }

  allocation_id = aws_eip.ipam[each.key].id
  subnet_id     = aws_subnet.subnet[each.key].id
}

resource "aws_vpc_endpoint" "endpoints" {
  for_each            = local.vpc_endpoints
  vpc_id              = aws_vpc.eks_vpc.id
  service_name        = each.value
  vpc_endpoint_type   = "Interface"
  security_group_ids  = [var.vpc_endpoints_security_group_id] # Ensure this allows HTTPS (443)
  subnet_ids          = [for s in keys(local.private_subnets) : aws_subnet.subnet[s].id]
  private_dns_enabled = true

  tags = {
    Name = "vpc-endpoint-${each.key}"
  }
}