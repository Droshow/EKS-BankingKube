module "networking" {
  source                        = "./modules/networking"
  alb_security_group            = [module.security.alb_sg_id]
  acm_domain_validation_options = module.security.domain_validation_options
  acm_certificate_arn           = module.security.certificate_arn
}

module "security" {
  source                  = "./modules/security"
  cluster_name            = var.cluster_name
  vpc_id                  = module.networking.vpc_id
  domain_name             = var.domain_name
  route_53cert_validation = module.networking.aws_route_53_cert_validation
  create_acm_certificate  = false
  tags = {
    Environment = "Banking-Kube"
  }
}