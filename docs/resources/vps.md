---
page_title: "liquidservers_vps Resource"
description: |-
  Manages a LiquidServers VPS.
---

# liquidservers_vps Resource

Creates, reads, upgrades, and terminates a VPS through the LiquidServers Server Portal API.

## Example Usage

```terraform
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

## Group Example

Every server must have a unique `client_reference`. The provider derives a per-server idempotency key from that reference.

```terraform
resource "liquidservers_vps" "workers" {
  for_each = var.servers

  hostname         = each.value.hostname
  label            = each.value.name
  client_reference = "terraform.workers.${each.key}"
  os_template      = each.value.os_template

  ipv4_count   = each.value.ipv4_count
  ram_mb       = each.value.ram_mb
  disk_gb      = each.value.disk_gb
  cpu_cores    = each.value.cpu_cores
  bandwidth_gb = var.default_bandwidth_gb
}
```

## Update Behavior

The portal supports resource increases through `/api/vps/upgrade`. Downgrades are rejected by the provider with a planning error. Changing hostname, label, OS template, node, root password, or client reference requires replacement.

## Import

Import by portal VM ID:

```shell
terraform import liquidservers_vps.web 57
```

## Schema

### Required

- `hostname` (String) Primary VPS hostname.
- `label` (String) Friendly label stored in the portal.
- `os_template` (String) Portal OS template slug.
- `ipv4_count` (Number) Number of IPv4 addresses to allocate to this VPS.
- `ram_mb` (Number) Memory in MB.
- `disk_gb` (Number) Disk allocation in GB.
- `cpu_cores` (Number) Number of CPU cores.
- `bandwidth_gb` (Number) Bandwidth allocation in GB.

### Optional

- `package_id` (Number) Optional portal package ID.
- `client_reference` (String) Stable per-server automation reference.
- `node_name` (String) Optional target node name.
- `location_name` (String) Optional location label.
- `root_password` (String, Sensitive) Optional initial root password.

### Read-Only

- `id` (String) Portal VM ID.
- `status` (String) Current portal status.
- `pending_action` (String) Current pending action.
- `ip_address` (String) Primary IP address when available.
- `provider_vps_id` (Number) Underlying provider VPS ID.
- `last_synced_at` (String) Last portal sync timestamp.

