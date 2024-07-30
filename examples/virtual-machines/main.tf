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
  node   = "pve"
  id     = 888
  cores  = 1
  memory = 512
}

output "cores" {
  value = proxmox_virtual_machine.vm1.cores
}

output "memory" {
  value = proxmox_virtual_machine.vm1.memory
}