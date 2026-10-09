terraform {

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "6.67.0"
    }
    github = {
      source  = "integrations/github"
      version = "6.13.0"
    }

  }
  required_version = "1.16.5"
}
