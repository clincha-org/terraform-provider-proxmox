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
  ide_devices = [
    {
      id = 1
      storage = "local-lvm"
      size = "5G"
    },
    {
      id = 2
      storage = "local"
      path = "iso/ubuntu-24.04.1-live-server-amd64.iso"
    }
  ]
}

output "cores" {
  value = proxmox_virtual_machine.vm1.cores
}

output "memory" {
  value = proxmox_virtual_machine.vm1.memory
}