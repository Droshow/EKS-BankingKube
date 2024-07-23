variable subnet_id {
  description = "The subnet ID to associate with the Client VPN"
  type        = string
}

variable server_certificate_arn {
  description = "The ARN of the server certificate"
  type        = string
}

variable client_root_certificate_arn {
  description = "The ARN of the client root certificate"
  type        = string
}