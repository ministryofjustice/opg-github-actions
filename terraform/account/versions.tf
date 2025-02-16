terraform {

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.86.1"
    }
    github = {
      source  = "integrations/github"
      version = "6.5.0"
    }

  }
  required_version = "1.10.5"
}
