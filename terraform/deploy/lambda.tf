resource "aws_lambda_function" "fn" {
  function_name    = "${var.fn}"
  filename         = "../../historian.zip"
  handler          = "historian"
  source_code_hash = "${base64sha256(file("../../historian.zip"))}"
  role             = "${aws_iam_role.fn.arn}"

  runtime     = "go1.x"
  memory_size = 128
  timeout     = 15

  environment {
    variables = {
      HISTORIAN_TABLE   = "${var.db}"
      HISTORIAN_USE_ALB = "${var.use_alb ? "true" : "false"}"
    }
  }
}

resource "aws_lambda_permission" "lb" {
  count = "${var.use_alb ? 1 : 0}"

  statement_id  = "AllowExecutionFromALB"
  action        = "lambda:InvokeFunction"
  function_name = "${aws_lambda_function.fn.arn}"
  principal     = "elasticloadbalancing.amazonaws.com"
  source_arn    = "${join("", aws_alb_target_group.historian.*.arn)}"
}

resource "aws_lambda_permission" "apigw" {
  count = "${var.use_alb ? 0 : 1}"

  statement_id  = "AllowExecutionFromAPIGW"
  action        = "lambda:InvokeFunction"
  function_name = "${aws_lambda_function.fn.arn}"
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${join("", aws_api_gateway_deployment.default.*.execution_arn)}*/*/*"
}
