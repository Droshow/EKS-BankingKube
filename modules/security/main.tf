resource "aws_security_group" "eks_cluster_sg" {
  name        = "${var.cluster_name}-sg"
  description = "Security group for the EKS cluster"
  vpc_id      = var.vpc_id
}

resource "aws_security_group" "worker_nodes_sg" {
  name        = "${var.cluster_name}-worker-nodes-sg"
  description = "Security group for the EKS worker nodes"
  vpc_id      = var.vpc_id
}

resource "aws_security_group" "alb_sg" {
  name        = "${var.cluster_name}-alb-sg"
  description = "Security group for the ALB"
  vpc_id      = var.vpc_id
}

# Allow the ALB to communicate with the worker nodes
resource "aws_security_group_rule" "alb_to_worker_nodes" {
  type                     = "ingress"
  from_port                = 0
  to_port                  = 0
  protocol                 = "-1"
  source_security_group_id = aws_security_group.alb_sg.id
  security_group_id        = aws_security_group.worker_nodes_sg.id
}

# Allow the worker nodes to communicate with the EKS cluster
resource "aws_security_group_rule" "worker_nodes_to_eks_cluster" {
  type                     = "ingress"
  from_port                = 0
  to_port                  = 0
  protocol                 = "-1"
  source_security_group_id = aws_security_group.worker_nodes_sg.id
  security_group_id        = aws_security_group.eks_cluster_sg.id
}