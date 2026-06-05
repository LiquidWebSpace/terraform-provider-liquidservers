---
page_title: "liquidservers_vps Data Source"
description: |-
  Reads an existing LiquidServers VPS.
---

# liquidservers_vps Data Source

Reads an existing VPS by portal VM ID.

## Example Usage

```terraform
data "liquidservers_vps" "existing" {
  id = "57"
}

output "existing_ip" {
  value = data.liquidservers_vps.existing.ip_address
}
```

## Schema

### Required

- `id` (String) Portal VM ID.

### Read-Only

- `hostname` (String)
- `label` (String)
- `package_id` (Number)
- `client_reference` (String)
- `os_template` (String)
- `node_name` (String)
- `location_name` (String)
- `ipv4_count` (Number)
- `ram_mb` (Number)
- `disk_gb` (Number)
- `cpu_cores` (Number)
- `bandwidth_gb` (Number)
- `status` (String)
- `pending_action` (String)
- `ip_address` (String)
- `provider_vps_id` (Number)
- `last_synced_at` (String)

