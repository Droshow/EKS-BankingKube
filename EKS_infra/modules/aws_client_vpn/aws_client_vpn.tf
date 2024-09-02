resource "aws_ec2_client_vpn_endpoint" "aws_client_vpn_endpoint" {
  description            = "AWS Client VPN Endpoint for secure access"
  server_certificate_arn = var.server_certificate_arn
  client_cidr_block      = "10.1.0.0/16"

  authentication_options {
    type                       = "certificate-authentication"
    root_certificate_chain_arn = var.client_root_certificate_arn
  }

  #  authentication_options {
  #   type                       = "mutual-authentication"
  #   root_certificate_chain_arn = var.client_root_certificate_arn
  # }

  connection_log_options {
    enabled               = true
    cloudwatch_log_group  = aws_cloudwatch_log_group.aws_client_vpn_logs.name
    cloudwatch_log_stream = aws_cloudwatch_log_stream.aws_client_vpn_log_stream.name
  }

  split_tunnel = true
  
  # Enable the self-service portal
  self_service_portal = "enabled"
}

resource "aws_cloudwatch_log_group" "aws_client_vpn_logs" {
  name = "aws-client-vpn-logs"
}

resource "aws_cloudwatch_log_stream" "aws_client_vpn_log_stream" {
  name           = "aws-client-vpn-log-stream"
  log_group_name = aws_cloudwatch_log_group.aws_client_vpn_logs.name
}

resource "aws_ec2_client_vpn_network_association" "aws_client_vpn_network_association" {
  client_vpn_endpoint_id = aws_ec2_client_vpn_endpoint.aws_client_vpn_endpoint.id
  subnet_id              = var.subnet_id
}

resource "aws_ec2_client_vpn_authorization_rule" "aws_client_vpn_authorization_rule" {
  client_vpn_endpoint_id = aws_ec2_client_vpn_endpoint.aws_client_vpn_endpoint.id
  target_network_cidr    = "0.0.0.0/0"
  authorize_all_groups   = true
}