# Terraform Provider for LiquidServers

This provider manages LiquidServers VPS resources through the LiquidServers Server Portal API.

It is structured for publishing to the Terraform Registry:

- Repository name: `terraform-provider-liquidservers`
- Provider address: `registry.terraform.io/liquidservers/liquidservers`
- Registry manifest: `terraform-registry-manifest.json`
- Registry documentation: `docs/`
- Release automation: `.goreleaser.yml` and `.github/workflows/release.yml`

## Example

```hcl
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

resource "liquidservers_vps" "web" {
  hostname         = "web-01.example.com"
  label            = "Web 01"
  client_reference = "terraform.production.web-01"
  os_template      = "ubuntu-24.04-x86_64"

  ipv4_count   = 1
  ram_mb       = 4096
  disk_gb      = 80
  cpu_cores    = 2
  bandwidth_gb = 3000
}
```

For groups of servers, set a unique `client_reference` per resource instance.

