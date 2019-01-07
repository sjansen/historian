resource "aws_iam_policy" "fn" {
  name = "${var.fn}"
  path = "/"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Resource": "arn:aws:logs:*:*:*"
    }, {
      "Effect": "Allow",
      "Action": [
        "dynamodb:PutItem"
      ],
      "Resource": "${data.aws_dynamodb_table.db.arn}"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "fn" {
  role = "${aws_iam_role.fn.name}"
  policy_arn = "${aws_iam_policy.fn.arn}"
}

resource "aws_iam_role" "fn" {
  name = "${var.fn}"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      }
    }
  ]
}
EOF
}
