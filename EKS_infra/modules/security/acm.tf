resource "aws_acm_certificate" "cert" {
  for_each          = var.certificate_names
  domain_name       = each.value
  validation_method = "DNS"
  tags              = lookup(var.tags, each.key, {})

  lifecycle {
    create_before_destroy = true
  }
}

#if this makes problems, then validation by hand in AWS is acceptable for now totally to simplify
resource "aws_acm_certificate_validation" "cert" {
  for_each                = aws_acm_certificate.cert
  certificate_arn         = aws_acm_certificate.cert[each.key].arn
  validation_record_fqdns = [for record in var.route_53cert_validation : record.fqdn]
}

### if cert already exists, then use this

# data "aws_acm_certificate" "existing_cert" {
#   for_each    = var.certificate_names
#   domain      = each.value
#   most_recent = true
#   statuses    = ["ISSUED"]
# }

