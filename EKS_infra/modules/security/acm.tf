resource "aws_acm_certificate" "cert" {
  for_each          = var.certificate_names
  domain_name       = each.value
  validation_method = "DNS"
  tags              = lookup(var.tags, each.key, {})

  lifecycle {
    prevent_destroy = true
  }
}

#if this makes problems, then validation by hand in AWS is acceptable for now totally to simplify
resource "aws_acm_certificate_validation" "cert" {
  for_each = {
    for key, cert in aws_acm_certificate.cert :
    key => cert if !var.fetch_existing_certificates || (
      var.fetch_existing_certificates && try(data.aws_acm_certificate.existing_cert[key].status, "NOT_ISSUED") != "ISSUED"
    )
  }
  certificate_arn         = each.value.arn
  validation_record_fqdns = try([for record in var.route_53cert_validation : record.fqdn], [])
}


### if cert already exists, then use this

# Fetch the existing certificate details
data "aws_acm_certificate" "existing_cert" {
  for_each    = var.fetch_existing_certificates ? var.certificate_names : {}
  domain      = each.value
  statuses    = ["ISSUED"]
  most_recent = true
}

