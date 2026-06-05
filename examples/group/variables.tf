variable "liquidservers_endpoint" {
  type    = string
  default = "https://servers.mywebsitepanel.com"
}

variable "liquidservers_api_key" {
  type      = string
  sensitive = true
}

variable "package_id" {
  type    = number
  default = null
}

variable "node_name" {
  type    = string
  default = null
}

variable "default_bandwidth_gb" {
  type    = number
  default = 3000
}

variable "servers" {
  type = map(object({
    name        = string
    hostname    = string
    os_template = string
    ipv4_count  = number
    ram_mb      = number
    disk_gb     = number
    cpu_cores   = number
  }))
}

