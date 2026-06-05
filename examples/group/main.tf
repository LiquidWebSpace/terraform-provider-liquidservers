terraform {
  required_providers {
    liquidservers = {
      source  = "liquidservers/liquidservers"
      version = "~> 0.1"
    }
  }
}

provider "liquidservers" {
  endpoint = var.liquidservers_endpoint
  api_key  = var.liquidservers_api_key
}

resource "liquidservers_vps" "servers" {
  for_each = var.servers

  hostname         = each.value.hostname
  label            = each.value.name
  client_reference = "terraform.group.${each.key}"
  os_template      = each.value.os_template

  package_id   = var.package_id
  node_name    = var.node_name
  ipv4_count   = each.value.ipv4_count
  ram_mb       = each.value.ram_mb
  disk_gb      = each.value.disk_gb
  cpu_cores    = each.value.cpu_cores
  bandwidth_gb = var.default_bandwidth_gb
}

output "server_ids" {
  value = { for key, server in liquidservers_vps.servers : key => server.id }
}

output "server_ips" {
  value = { for key, server in liquidservers_vps.servers : key => server.ip_address }
}

