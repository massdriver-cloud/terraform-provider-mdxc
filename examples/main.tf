terraform {
  required_providers {
    massdriver = {
      version = "0.0.1"
      source  = "massdriver.cloud/mdxc"
    }
  }
}

provider "massdriver" {}
