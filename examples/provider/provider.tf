terraform {
  required_providers {
    liquidservers = {
      source  = "liquidservers/liquidservers"
      version = "~> 0.1"
    }
  }
}

provider "liquidservers" {
  endpoint = "https://servers.mywebsitepanel.com"
  api_key  = var.liquidservers_api_key
}

variable "liquidservers_api_key" {
  type      = string
  sensitive = true
}

