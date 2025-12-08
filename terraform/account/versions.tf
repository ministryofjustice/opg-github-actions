terraform {

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "6.20.0"
    }
    github = {
      source  = "integrations/github"
      version = "6.7.5"
    }

  }
  required_version = "1.14.1"
}
