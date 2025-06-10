terraform {
  required_providers {
    unifi = {
      source = "zoullx/unifi"
    }
  }
}

provider "unifi" {
  host           = "https://10.0.0.1"
  api_key        = "KjNe0gwYk1LPh-odK_9fkptG6_S5xaZV"
  allow_insecure = true
}

import {
  to = unifi_network.test
  id = "default/5cddfa45a4ccf30ec588017a"
}

resource "unifi_network" "test" {
  name               = "Network"
  site               = "default"
  setting_preference = "manual"
  dhcp_v6_dns_auto   = true
  ipv6_pd_stop       = "::7d1"
  # dhcp_gateway_enabled           = false
  ipv6_client_address_assignment = "slaac"
  dhcp_start                     = "10.0.0.2"
  ipv6_ra_enabled                = true
  domain_name                    = "zdl.io"
  subnet                         = "10.0.0.1/24"
  ipv6_interface_type            = "none"
  dhcp_v6_stop                   = "::7d1"
  dhcp_dns_enabled               = true
  # dhcp_v6_enabled                = false
  internet_access_enabled = true
  # dhcp_relay_enabled             = false
  dhcp_conflict_checking = true
  # dhcp_wins_enabled              = false
  ipv6_pd_auto_prefixid_enabled = true
  dhcp_v6_lease_time            = 86400
  lte_lan_enabled               = true
  dhcp_lease_time               = 86400
  purpose                       = "corporate"
  # igmp_snooping                  = false
  # dhcp_time_offset_enabled       = false
  # dhcp_guard_enabled             = false
  dhcp_v6_allow_slaac        = true
  ipv6_ra_preferred_lifetime = 14400
  dhcp_stop                  = "10.0.0.254"
  enabled                    = true
  dhcp_enabled               = true
  vlan_id                    = 0
  network_group              = "LAN"
  dhcp_v6_start              = "::2"
  # vlan_enabled                   = false
  ipv6_setting_preference = "auto"
  gateway_type            = "default"
  ipv6_ra_priority        = "high"
  # dhcp_boot_enabled              = false
  ipv6_pd_start = "::2"
  # upnp_lan_enabled               = false
  # dhcp_ntp_enabled               = false
  multicast_dns_enabled = true
  # auto_scale_enabled             = false
}

output "network_test" {
  value = unifi_network.test
}
