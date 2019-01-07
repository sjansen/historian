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
      HISTORIAN_TABLE = "${var.db}"
    }
  }
}

resource "aws_lambda_permission" "lb" {
  statement_id  = "AllowExecutionFromLB"
  action        = "lambda:InvokeFunction"
  function_name = "${aws_lambda_function.fn.arn}"
  principal     = "elasticloadbalancing.amazonaws.com"
  source_arn    = "${aws_alb_target_group.historian.arn}"
}
