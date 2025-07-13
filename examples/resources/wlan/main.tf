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

resource "unifi_wlan" "test" {
  name                                 = "Network"
  site                                 = "default"
  setting_preference                   = "manual"
  dtim_6e                              = 3
  wpa_mode                             = "wpa2"
  minimum_data_rate_setting_preference = "auto"
  radius_mac_acl_format                = "none_lower"
  pmf_mode                             = "disabled"
  user_group_id                        = "5cddfa45a4ccf30ec588017b"
  iapp_enabled                         = true
  wlan_band                            = "both"
  network_id                           = "6872c6d40a9c2e5995e4c068"
  dtim_5g                              = 3
  enabled                              = true
  wlan_bands                           = ["2g", "5g"]
  mac_filter_policy                    = "allow"
  security                             = "wpapsk"
  ap_group_ids                         = ["602cbd66a4ccf3212df26775"]
  minimum_2g_data_rate_enabled         = true
  bss_transition                       = true
  minimum_2g_data_rate_kbps            = 1000
  ap_group_mode                        = "all"
  wpa_enc                              = "ccmp"
  passphrase                           = "supersecretpassword"
  dtim_mode                            = "default"
  minimum_5g_data_rate_kbps            = 6000
  dtim_2g                              = 1
}

output "wlan_test" {
  value = unifi_wlan.test
}
