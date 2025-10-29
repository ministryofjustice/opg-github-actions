terraform {

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "6.18.0"
    }
    github = {
      source  = "integrations/github"
      version = "6.7.2"
    }

  }
  required_version = "1.13.4"
}
