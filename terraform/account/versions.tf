terraform {

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "6.42.0"
    }
    github = {
      source  = "integrations/github"
      version = "6.12.0"
    }

  }
  required_version = "1.14.9"
}
