terraform {
  required_providers {
    massdriver = {
      version = "~> 1.0.0"
      source  = "massdriver-cloud/mdxc"
    }
  }
}

variable "md_name_prefix" {
  type    = string
  default = "project-target-network-1234"
}


locals {
  aws_vpc = {
    main = {
      "arn" = "some fake arn"
      "id"  = "some fake id"
    }
  }

  aws_iam_role = {
    "foo" = {
      "arn" = "some fake arn"
    }
  }
}

provider "massdriver" {}
