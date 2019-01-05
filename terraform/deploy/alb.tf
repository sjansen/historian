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

resource "aws_alb_listener" "lb" {
  load_balancer_arn = "${aws_alb.lb.arn}"
  port              = "80"
  protocol          = "HTTP"

  default_action {
    target_group_arn = "${aws_alb_target_group.historian.arn}"
    type             = "forward"
  }
}

resource "aws_alb_target_group" "historian" {
  name        = "historian"
  target_type = "lambda"

  health_check {    
    healthy_threshold   = 2
    unhealthy_threshold = 5    
    timeout             = 5    
    interval            = 120    
    matcher             = "200"
    path                = "/health/"    
    port                = "traffic-port"
    protocol            = "HTTP"
  }
}

resource "aws_alb_target_group_attachment" "fn" {
  target_group_arn  = "${aws_alb_target_group.historian.arn}"
  target_id         = "${aws_lambda_function.fn.arn}"
  depends_on        = ["aws_lambda_permission.lb"]
}
