resource "aws_route53_record" "cert_validation" {
  for_each = {
    for dvo in var.acm_domain_validation_options : dvo.domain_name => {
      name   = dvo.resource_record_name
      record = dvo.resource_record_value
      type   = dvo.resource_record_type
    }
  }

  name    = each.value.name
  type    = each.value.type
  zone_id = aws_route53_zone.banking-kube.id
  records = [each.value.record]
  ttl     = 60
}

resource "aws_route53_zone" "banking-kube" {
  name = "bankingkube.com"

  comment = "Public DNS zone for Banking-Kube"
  tags = {
    Environment = "Banking-Kube"
  }

}