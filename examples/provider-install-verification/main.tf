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
