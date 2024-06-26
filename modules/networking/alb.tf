locals {
  listeners = {
    http = {
      load_balancer_arn = aws_lb.eks_alb.arn
      port              = 80
      protocol          = "HTTP"
      target_group_arn  = aws_lb_target_group.eks_tg.arn
    },
    https = {
      load_balancer_arn = aws_lb.eks_alb.arn
      port              = 443
      protocol          = "HTTPS"
      ssl_policy        = "ELBSecurityPolicy-2016-08"
      certificate_arn   = var.acm_certificate_arn
      target_group_arn  = aws_lb_target_group.eks_tg.arn
    }
  }
}
resource "aws_lb" "eks_alb" {
  name               = "eks-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = var.alb_security_group
  subnets            = [for _, subnet in aws_subnet.subnet : subnet.id if subnet.map_public_ip_on_launch]

  enable_deletion_protection = true

  tags = {
    Environment = "production"
  }
}
resource "aws_lb_listener" "eks_listener" {
  for_each = local.listeners

  load_balancer_arn = each.value.load_balancer_arn
  port              = each.value.port
  protocol          = each.value.protocol

  default_action {
    type             = "forward"
    target_group_arn = each.value.target_group_arn
  }

  certificate_arn = each.value.protocol == "HTTPS" ? each.value.certificate_arn : null
  ssl_policy      = each.value.protocol == "HTTPS" ? each.value.ssl_policy : null
}

resource "aws_lb_target_group" "eks_tg" {
  name     = "eks-tg"
  port     = 80
  protocol = "HTTP"
  vpc_id   = aws_vpc.eks_vpc.id

  health_check {
    enabled             = true
    interval            = 30
    path                = "/"
    port                = "traffic-port"
    timeout             = 3
    healthy_threshold   = 3
    unhealthy_threshold = 3
    matcher             = "200-399"
  }
}

