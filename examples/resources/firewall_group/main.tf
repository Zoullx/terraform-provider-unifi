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

resource "unifi_firewall_group" "test" {
  name    = "All Local"
  site    = "default"
  type    = "address-group"
  members = ["10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"]
}

output "firewall_group_test" {
  value = unifi_firewall_group.test
}
