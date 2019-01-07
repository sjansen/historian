data "aws_dynamodb_table" "db" {
  name         = "${var.db}"
}
