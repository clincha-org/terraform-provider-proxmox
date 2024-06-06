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

resource "proxmox_network" "vmbr88" {
  interface = "vmbr22"
  type      = "bridge"
  address = "10.2.0.34"
  netmask = "255.255.255.0"
  autostart = false
}

output "vmbr88" {
  value = proxmox_network.vmbr88
}