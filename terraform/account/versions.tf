terraform {

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.66.0"
    }
    github = {
      source  = "integrations/github"
      version = "5.45.0"
    }

  }
  required_version = "1.9.5"
}
