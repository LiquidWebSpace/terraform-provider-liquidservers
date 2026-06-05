---
page_title: "liquidservers_os_templates Data Source"
description: |-
  Lists available LiquidServers OS templates.
---

# liquidservers_os_templates Data Source

Lists OS templates available from the LiquidServers Server Portal API.

## Example Usage

```terraform
data "liquidservers_os_templates" "available" {}

output "template_slugs" {
  value = [for template in data.liquidservers_os_templates.available.templates : template.slug]
}
```

## Schema

### Read-Only

- `templates` (List of Object)
  - `id` (Number) Template ID.
  - `name` (String) Template display name.
  - `slug` (String) Template slug for VPS creation.
  - `osid` (Number) Provider OS ID.
  - `virt` (String) Virtualization type.

