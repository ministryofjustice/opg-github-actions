terraform {

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "6.42.0"
    }
    github = {
      source  = "integrations/github"
      version = "6.12.1"
    }

  }
  required_version = "1.15.0"
}
