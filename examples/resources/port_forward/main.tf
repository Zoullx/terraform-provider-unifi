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

resource "unifi_port_forward" "test" {
  name                   = "HTTP"
  site                   = "default"
  port_forward_interface = "wan"
  src_ip                 = "any"
  dst_port               = "80"
  fwd_ip                 = "10.0.1.2"
  fwd_port               = "80"
  protocol               = "tcp_udp"
  enabled                = true
}

output "port_forward_test" {
  value = unifi_port_forward.test
}
