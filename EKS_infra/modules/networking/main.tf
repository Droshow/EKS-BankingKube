locals {
  private_subnets = { for k, v in var.subnets : k => v if !v.public }
  public_subnets  = { for k, v in var.subnets : k => v if v.public }
}


resource "aws_vpc" "eks_vpc" {
  cidr_block = var.vpc_cidr_block
  enable_dns_support = true
}
resource "aws_subnet" "subnet" {
  for_each = var.subnets

  vpc_id                  = aws_vpc.eks_vpc.id
  cidr_block              = each.value.cidr
  map_public_ip_on_launch = each.value.public
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