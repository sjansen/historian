resource "aws_api_gateway_rest_api" "gw" {
  count = "${var.use_alb ? 0 : 1}"

  name        = "${var.lb}"
  endpoint_configuration {
    types = ["REGIONAL"]
  }
}

resource "aws_api_gateway_base_path_mapping" "default" {
  count = "${var.use_alb ? 0 : 1}"

  api_id      = "${join("", aws_api_gateway_rest_api.gw.*.id)}"
  stage_name  = "${join("", aws_api_gateway_deployment.default.*.stage_name)}"
  domain_name = "${join("", aws_api_gateway_domain_name.gw.*.domain_name)}"
}

resource "aws_api_gateway_deployment" "default" {
  count = "${var.use_alb ? 0 : 1}"

  depends_on = [
    "aws_api_gateway_integration.lambda",
    "aws_api_gateway_integration.lambda_root",
  ]

  rest_api_id = "${join("", aws_api_gateway_rest_api.gw.*.id)}"
  stage_name  = "default"
}

resource "aws_api_gateway_domain_name" "gw" {
  count = "${var.use_alb ? 0 : 1}"

  domain_name              = "${var.dns_name}"
  regional_certificate_arn = "${aws_acm_certificate_validation.cert.certificate_arn}"
  endpoint_configuration {
    types = ["REGIONAL"]
  }
}

resource "aws_api_gateway_method" "proxy" {
  count = "${var.use_alb ? 0 : 1}"

  rest_api_id   = "${join("", aws_api_gateway_rest_api.gw.*.id)}"
  resource_id   = "${join("", aws_api_gateway_resource.proxy.*.id)}"
  http_method   = "ANY"
  authorization = "NONE"
}

resource "aws_api_gateway_method" "proxy_root" {
  count = "${var.use_alb ? 0 : 1}"

  rest_api_id = "${join("", aws_api_gateway_rest_api.gw.*.id)}"
  resource_id = "${join("", aws_api_gateway_rest_api.gw.*.root_resource_id)}"
  http_method   = "ANY"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "lambda" {
  count = "${var.use_alb ? 0 : 1}"

  rest_api_id = "${join("", aws_api_gateway_rest_api.gw.*.id)}"
  resource_id = "${join("", aws_api_gateway_method.proxy.*.resource_id)}"
  http_method = "${join("", aws_api_gateway_method.proxy.*.http_method)}"

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = "${aws_lambda_function.fn.invoke_arn}"
}

resource "aws_api_gateway_integration" "lambda_root" {
  count = "${var.use_alb ? 0 : 1}"

  rest_api_id = "${join("", aws_api_gateway_rest_api.gw.*.id)}"
  resource_id = "${join("", aws_api_gateway_method.proxy_root.*.resource_id)}"
  http_method = "${join("", aws_api_gateway_method.proxy_root.*.http_method)}"

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = "${aws_lambda_function.fn.invoke_arn}"
}

resource "aws_api_gateway_resource" "proxy" {
  count = "${var.use_alb ? 0 : 1}"

  rest_api_id = "${join("", aws_api_gateway_rest_api.gw.*.id)}"
  parent_id   = "${join("", aws_api_gateway_rest_api.gw.*.root_resource_id)}"
  path_part   = "{proxy+}"
}
