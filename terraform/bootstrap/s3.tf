data "aws_elb_service_account" "main" {}

resource "aws_s3_bucket" "logs" {
  bucket        = "${var.logs}"
  acl           = "private"
  force_destroy = true
  lifecycle_rule {
    id      = "cleanup"
    enabled = true
    abort_incomplete_multipart_upload_days = 3
    expiration {
      days = 90
    }
    noncurrent_version_expiration {
      days = 30
    }
  }
  policy = <<POLICY
{
  "Id": "Policy",
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "s3:PutObject"
      ],
      "Effect": "Allow",
      "Resource": "arn:aws:s3:::${var.logs}/AWSLogs/*",
      "Principal": {
        "AWS": [
          "${data.aws_elb_service_account.main.arn}"
        ]
      }
    }
  ]
}
POLICY
  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        sse_algorithm = "AES256"
      }
    }
  }
  versioning {
    enabled = "${var.protect_logs}"
  }
}
