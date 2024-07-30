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

resource "proxmox_virtual_machine" "vm1" {
  node = "pve"
  id   = 888
}
#
# output "cores" {
#   value = vm1.cores
# }
#
# output "memory" {
#   value = vm1.memory
# }