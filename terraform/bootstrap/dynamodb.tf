resource "aws_dynamodb_table" "db" {
  name         = "${var.db}"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "event-id"
  range_key    = "timestamp"

  attribute {
    name = "event-id"
    type = "S"
  }

  attribute {
    name = "timestamp"
    type = "N"
  }

  server_side_encryption {
    enabled = true
  }

  ttl {
    attribute_name = "expires"
    enabled        = true
  }
}
