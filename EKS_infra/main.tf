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
    module.security.eks_cluster_sg_id,
    module.security.worker_nodes_sg_id,
    module.security.alb_sg_id
  ]
  cluster_role_iam_role_arn = module.security.eks_cluster_role_arn
}

module "node_groups" {
  source                     = "./modules/node_groups"
  cluster_name               = var.cluster_name
  fargate_pod_execution_role = module.security.fargate_pod_execution_role_arn
}

module "db_instance" {
  source         = "./modules/databases"
  db_name        = "banking-kube-db"
  username       = "Dro-admin"
  password       = random_password.db_password.result
  db_subnets     = module.networking.private_subnets_ids
}

module "storage" {
  source     = "./modules/storage"
  subnet_ids = module.networking.private_subnets_ids
  vpc_id     = module.networking.vpc_id
  tags = {
    Environment = "Banking-Kube"
  }
  efs_sg_id = module.security.efs_security_group_id
}