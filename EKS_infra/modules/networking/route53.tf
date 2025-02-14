resource "aws_route53_record" "cert_validation" {
  for_each = {
    for cert_name, validation_options in var.acm_domain_validation_options : cert_name => validation_options
  }

  name    = each.value[0].resource_record_name
  type    = each.value[0].resource_record_type
  zone_id = aws_route53_zone.banking-kube.id
  records = [each.value[0].resource_record_value]
  ttl     = 60
}


resource "aws_route53_zone" "banking-kube" {
  name = "devsbridge.com"

  comment = "Public DNS zone for Banking-Kube"
  tags = {
    Environment = "DevsBridge"
  }
  lifecycle {
    prevent_destroy = false
  }

}