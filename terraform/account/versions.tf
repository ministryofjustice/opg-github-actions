terraform {

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "6.25.0"
    }
    github = {
      source  = "integrations/github"
      version = "6.10.2"
    }

  }
  required_version = "1.14.3"
}
