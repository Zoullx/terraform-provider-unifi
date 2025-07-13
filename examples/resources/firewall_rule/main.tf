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

resource "unifi_firewall_rule" "test" {
  name               = "Allow Established and Related to Any"
  site               = "default"
  setting_preference = "manual"
  src_network_type   = "NETv4"
  state_established  = true
  enabled            = true
  protocol           = "all"
  action             = "accept"
  dst_network_type   = "NETv4"
  state_related      = true
  ruleset            = "LAN_IN"
  rule_index         = 20000
}

output "firewall_rule_test" {
  value = unifi_firewall_rule.test
}
