resource "liquidservers_vps" "web" {
  hostname         = "web-01.example.com"
  label            = "Web 01"
  client_reference = "terraform.example.web-01"
  os_template      = "ubuntu-24.04-x86_64"

  ipv4_count   = 1
  ram_mb       = 4096
  disk_gb      = 80
  cpu_cores    = 2
  bandwidth_gb = 3000
}

