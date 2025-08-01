terraform {

  backend "s3" {
    bucket  = "opg.terraform.state"
    key     = "github-actions/account/terraform.tfstate"
    encrypt = true
    region  = "eu-west-1"
    assume_role = {
      role_arn = "arn:aws:iam::311462405659:role/gh-reusable-actions-ci"
    }
    dynamodb_table = "remote_lock"
  }

}

provider "github" {
  token = var.github_token
  owner = "ministryofjustice"
}

locals {
  sandbox = "995199299616"
}

provider "aws" {
  alias = "sandbox"

  assume_role {
    role_arn     = "arn:aws:iam::${local.sandbox}:role/${var.DEFAULT_ROLE}"
    session_name = "terraform-session"
  }
}

provider "aws" {
  assume_role {
    role_arn     = "arn:aws:iam::${local.sandbox}:role/${var.DEFAULT_ROLE}"
    session_name = "terraform-session"
  }
}

variable "github_token" {
}

variable "aws_access_key_id" {
}

variable "aws_secret_access_key" {
}

variable "DEFAULT_ROLE" {
  default = "gh-reusable-actions-ci"
}

data "aws_caller_identity" "current" {
  provider = aws.sandbox
}
