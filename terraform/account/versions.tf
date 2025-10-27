terraform {

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "6.17.0"
    }
    github = {
      source  = "integrations/github"
      version = "6.7.0"
    }

  }
  required_version = "1.13.4"
}
