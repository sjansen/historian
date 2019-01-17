data "aws_route53_zone" "zone" {
  name = "${var.dns_zone}"
  private_zone = false
}


resource "aws_route53_record" "cert" {
  zone_id = "${data.aws_route53_zone.zone.id}"
  name = "${aws_acm_certificate.cert.domain_validation_options.0.resource_record_name}"
  type = "${aws_acm_certificate.cert.domain_validation_options.0.resource_record_type}"
  ttl = 60
  records = [
    "${aws_acm_certificate.cert.domain_validation_options.0.resource_record_value}"
  ]
}


resource "aws_route53_record" "alb" {
  count = "${var.use_alb ? 1 : 0}"

  zone_id  = "${data.aws_route53_zone.zone.id}"
  name     = "${var.dns_name}"
  type     = "A"
  alias {
    name     = "${join("", aws_alb.lb.*.dns_name)}"
    zone_id  = "${join("", aws_alb.lb.*.zone_id)}"
    evaluate_target_health = false
  }
}

resource "aws_route53_record" "apigw" {
  count = "${var.use_alb ? 0 : 1}"

  zone_id  = "${data.aws_route53_zone.zone.id}"
  name    = "${join("", aws_api_gateway_domain_name.gw.*.domain_name)}"
  type    = "A"
  alias {
    name     = "${join("", aws_api_gateway_domain_name.gw.*.regional_domain_name)}"
    zone_id  = "${join("", aws_api_gateway_domain_name.gw.*.regional_zone_id)}"
    evaluate_target_health = false
  }
}
