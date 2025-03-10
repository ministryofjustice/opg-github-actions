terraform {

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.90.0"
    }
    github = {
      source  = "integrations/github"
      version = "6.6.0"
    }

  }
  required_version = "1.11.1"
}
