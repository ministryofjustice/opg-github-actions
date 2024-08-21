terraform {

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.63.0"
    }
    github = {
      source  = "integrations/github"
      version = "5.45.0"
    }

  }
  required_version = "1.9.4"
}
