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

variable "dns_name" {
  type = "string"
}

variable "dns_zone" {
  type = "string"
}

variable "fn" {
  type = "string"
}

variable "publish_fn" {
  default = false
}

variable "lb" {
  type = "string"
}

variable "protect_lb" {
  default = false
}

variable "logs" {
  type = "string"
}

variable "sg" {
  type = "string"
}

variable "subnet_ids" {
  type = "list"
}

variable "vpc_id" {
  type = "string"
}
