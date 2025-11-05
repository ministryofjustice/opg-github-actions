terraform {

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "6.19.0"
    }
    github = {
      source  = "integrations/github"
      version = "6.7.4"
    }

  }
  required_version = "1.13.4"
}
