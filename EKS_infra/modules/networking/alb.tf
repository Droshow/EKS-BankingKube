locals {
  listeners = {
    http = {
      load_balancer_arn = aws_lb.eks_alb.arn
      port              = 80
      protocol          = "HTTP"
      target_group_arn  = aws_lb_target_group.eks_tg.arn
    },
   # TODO: Uncomment when certificate is created
    https = {
      load_balancer_arn = aws_lb.eks_alb.arn
      port              = 443
      protocol          = "HTTPS"
      ssl_policy        = "ELBSecurityPolicy-2016-08"
      certificate_arn   = var.acm_certificate_arn
      target_group_arn  = aws_lb_target_group.eks_tg.arn
    }
  }
  #   public_subnets_alb = [aws_subnet.subnet["eks_public_subnet_001"].id, aws_subnet.subnet["eks_public_subnet_002"].id
  # ]
  # public_subnet = var.subnets["eks_public_subnet-001"].public ? aws_subnet.subnet["eks_public_subnet-001"].id : aws_subnet.subnet["eks_public_subnet-002"].id
  public_subnets_per_az = { for k, v in var.subnets : v.az => aws_subnet.subnet[k].id if v.public }
}
resource "aws_lb" "eks_alb" {
  name               = "eks-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = var.alb_security_group
  subnets            = values(local.public_subnets_per_az)

  #TODO usually true in production areas for test false
  enable_deletion_protection = false

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

