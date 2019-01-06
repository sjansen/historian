resource "aws_alb" "lb" {
  name               = "${var.lb}"
  internal           = false
  load_balancer_type = "application"
  security_groups    = ["${aws_security_group.sg.id}"]
  subnets            = "${var.subnet_ids}"

  enable_deletion_protection = "${var.protect_lb}"

  access_logs {
    bucket  = "${data.aws_s3_bucket.logs.bucket}"
    enabled = true
  }
}

resource "aws_alb_listener" "http" {
  load_balancer_arn = "${aws_alb.lb.arn}"

  port     = "80"
  protocol = "HTTP"

  default_action {
    type = "redirect"
    redirect {
      host        = "${var.dns_name}"
      port        = "443"
      protocol    = "HTTPS"
      status_code = "HTTP_302"
    }
  }
}

resource "aws_alb_listener" "https" {
  certificate_arn   = "${aws_acm_certificate.cert.arn}"
  load_balancer_arn = "${aws_alb.lb.arn}"

  port       = "443"
  protocol   = "HTTPS"
  ssl_policy = "ELBSecurityPolicy-FS-2018-06"

  default_action {
    target_group_arn = "${aws_alb_target_group.historian.arn}"
    type             = "forward"
  }
}

resource "aws_alb_target_group" "historian" {
  name        = "historian"
  target_type = "lambda"
}

resource "aws_alb_target_group_attachment" "fn" {
  target_group_arn  = "${aws_alb_target_group.historian.arn}"
  target_id         = "${aws_lambda_function.fn.arn}"
  depends_on        = ["aws_lambda_permission.lb"]
}
