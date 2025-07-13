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

resource "unifi_network" "test" {
  name                           = "Servers"
  site                           = "default"
  setting_preference             = "manual"
  dhcp_v6_dns_auto               = true
  ipv6_pd_stop                   = "::7d1"
  ipv6_client_address_assignment = "slaac"
  dhcp_start                     = "10.1.0.2"
  ipv6_ra_enabled                = true
  domain_name                    = "zdl.io"
  subnet                         = "10.1.0.1/24"
  ipv6_interface_type            = "none"
  dhcp_v6_stop                   = "::7d1"
  internet_access_enabled        = true
  dhcp_relay_enabled             = true
  dhcp_conflict_checking         = true
  ipv6_pd_auto_prefixid_enabled  = true
  dhcp_v6_lease_time             = 86400
  lte_lan_enabled                = true
  dhcp_lease_time                = 86400
  purpose                        = "corporate"
  dhcp_v6_allow_slaac            = true
  ipv6_ra_preferred_lifetime     = 14400
  dhcp_stop                      = "10.1.0.254"
  enabled                        = true
  dhcp_enabled                   = false
  vlan_id                        = 2
  network_group                  = "LAN"
  dhcp_v6_start                  = "::2"
  vlan_enabled                   = true
  ipv6_setting_preference        = "manual"
  gateway_type                   = "default"
  ipv6_ra_priority               = "high"
  ipv6_pd_start                  = "::2"
  multicast_dns_enabled          = true
}

output "network_test" {
  value = unifi_network.test
}
