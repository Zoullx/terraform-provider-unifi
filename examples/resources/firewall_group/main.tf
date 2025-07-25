terraform {
  required_providers {
    unifi = {
      source = "zoullx/unifi"
    }
  }
}

variable "UNIFI_API_KEY" {
  type        = string
  description = "Unifi API Key"
  sensitive   = true
}

provider "unifi" {
  host           = "https://10.0.0.1"
  api_key        = var.UNIFI_API_KEY
  allow_insecure = true
}

resource "unifi_firewall_group" "test" {
  name    = "All Local"
  site    = "default"
  type    = "address-group"
  members = ["10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"]
}

output "firewall_group_test" {
  value = unifi_firewall_group.test
}
