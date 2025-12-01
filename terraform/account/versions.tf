terraform {

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "6.23.0"
    }
    github = {
      source  = "integrations/github"
      version = "6.8.3"
    }

  }
  required_version = "1.14.0"
}
