terraform {
  required_providers {
    proxmox = {
      source = "hashicorp.com/edu/proxmox"
    }
  }
}

provider "proxmox" {
  host     = "https://localhost:8006"
  username = "root@pam"
  password = "vagrant"
}

data "proxmox_node" "edu" {}

output "edu_nodes" {
  value = data.proxmox_node.edu
}

resource "proxmox_network" "edu" {
  iface     = "vmbr88"
  type      = "bridge"
  autostart = 1
}

resource "proxmox_network" "edu1" {
  iface = "vmbr3"
  type  = "bridge"
}

resource "proxmox_network" "edu5" {
  iface = "vmbr7"
  type  = "bridge"
}

output "edu_network" {
  value = proxmox_network.edu
}