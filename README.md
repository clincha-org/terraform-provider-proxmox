# Terraform Provider for Proxmox

## Provider Configuration

Only password authentication is supported at the moment. TLS certificate verification is disabled.

```hcl
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
```

## Usage

At the moment there is only one resource and one data source available.

### Resource `proxmox_network_bridge`

```hcl
resource "proxmox_network_bridge" "vmbr88" {
  interface = "vmbr88"
  node      = "pve"
  address   = "10.1.2.3"
  netmask   = "255.255.255.0"
  autostart = false
  comments  = "Terraform created network"
}
```

### Data Source `proxmox_node`

This data source returns information about all the proxmox **nodes** in the cluster. It returns a list of nodes. 

```hcl
data "proxmox_node" "edu" {

}
```
