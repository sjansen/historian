output "lb" {
  value = "${aws_alb.lb.dns_name}"
}
