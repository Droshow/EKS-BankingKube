#probably should serve as a constructor resource that can output whatever acm certificates are needed
resource "aws_acm_certificate" "cert" {
  count             = var.create_acm_certificate ? 1 : 0
  domain_name       = var.domain_name
  validation_method = "DNS"
  tags              = var.tags

  lifecycle {
    create_before_destroy = true
  }
}
#if this makes problems, then validation by hand in AWS is acceptable for now totally to simplify
resource "aws_acm_certificate_validation" "cert" {
  count                   = var.create_acm_certificate ? 1 : 0
  certificate_arn         = aws_acm_certificate.cert[0].arn
  validation_record_fqdns = [for record in var.route_53cert_validation : record.fqdn]
}

