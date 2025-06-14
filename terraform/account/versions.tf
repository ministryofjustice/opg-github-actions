terraform {

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.99.1"
    }
    github = {
      source  = "integrations/github"
      version = "6.6.0"
    }

  }
  required_version = "1.12.1"
}
