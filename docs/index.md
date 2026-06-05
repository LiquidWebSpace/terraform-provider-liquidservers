---
page_title: "LiquidServers Provider"
description: |-
  Manage LiquidServers VPS resources through the LiquidServers Server Portal API.
---

# LiquidServers Provider

The LiquidServers provider creates, reads, upgrades, and terminates VPS resources through the LiquidServers Server Portal API.

The portal remains the control plane: reseller ownership, package limits, API logging, provider payload shaping, and Virtualizor integration all stay inside the portal.

## Example Usage

```terraform
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
```

## Authentication

Use a portal API key generated from the LiquidServers API Help page.

The provider accepts credentials in the provider block or through environment variables:

- `LIQUIDSERVERS_ENDPOINT`
- `LIQUIDSERVERS_API_KEY`

## Publishing Notes

This repository is structured for Terraform Registry publishing:

- GitHub repository name must be `terraform-provider-liquidservers`
- The repository must be public and lowercase
- Releases must use semantic version tags such as `v0.1.0`
- Release assets must be signed with a GPG key registered in the Terraform Registry
- `terraform-registry-manifest.json` declares Terraform Plugin Framework protocol `6.0`

## Schema

### Optional

- `endpoint` (String) LiquidServers Server Portal base URL. Can also be set with `LIQUIDSERVERS_ENDPOINT`.
- `api_key` (String, Sensitive) LiquidServers API key. Can also be set with `LIQUIDSERVERS_API_KEY`.

