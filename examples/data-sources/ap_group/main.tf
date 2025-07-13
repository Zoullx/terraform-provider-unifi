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

data "unifi_ap_group" "test" {
  name = "All APs"
  site = "default"
}

output "ap_group_test" {
  value = data.unifi_ap_group.test
}
