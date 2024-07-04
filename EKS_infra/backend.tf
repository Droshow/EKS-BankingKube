terraform {
  backend "s3" {
    bucket = "banking-kube-state-tf"
    key    = "terraform.tfstate"
    region = "eu-west-1"
  }
}