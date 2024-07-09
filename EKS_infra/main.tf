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
module "eks" {
  source       = "./modules/eks"
  cluster_name = var.cluster_name
  subnet_ids   = module.networking.private_subnets_ids
  security_group_ids = [
    module.outputs.eks_cluster_sg_id,
    module.outputs.worker_nodes_sg_id,
    module.outputs.alb_sg_id
  ]
  cluster_role_iam_role_arn = module.security.eks_cluster_role_arn
}

module "node_groups" {
  source                     = "./modules/node_groups"
  cluster_name               = var.cluster_name
  fargate_pod_execution_role = module.security.fargate_pod_execution_role_arn


}
