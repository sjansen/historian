terraform {
  required_version = ">= 0.11.11"
}

##
# Credentials
##

provider "aws" {
  version = "~> 1.54"

  profile = "${var.aws_profile}"
  region = "${var.aws_region}"
}

variable "aws_profile" {
  type = "string"
}

variable "aws_region" {
  default = "us-east-1"
}

##
# Resources
##

variable "db" {
  type = "string"
}

variable "logs" {
  type = "string"
}

variable "protect_logs" {
  default = false
}
