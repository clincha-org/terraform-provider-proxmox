package provider

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestVirtualMachine_Create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "proxmox_virtual_machine" "vm1" {
  node   = "pve"
  id     = 888
  cores  = 1
  memory = 512
  ide_devices = [
	{
      id = 1
	  storage = "local-lvm"
      size = "8G"
     },
	{
      id = 2
	  storage = "local"
      path = "iso/ubuntu-24.04.1-live-server-amd64.iso"
     }
  ]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_virtual_machine.vm1", "cores", "1"),
					resource.TestCheckResourceAttr("proxmox_virtual_machine.vm1", "memory", "512"),
				),
			},
			{
				ResourceName:      "proxmox_virtual_machine.vm1",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "pve,888",
			},
			{
				Config: providerConfig + `
resource "proxmox_virtual_machine" "vm1" {
  node   = "pve"
  id     = 888
  cores  = 1
  memory = 1024
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_virtual_machine.vm1", "cores", "1"),
					resource.TestCheckResourceAttr("proxmox_virtual_machine.vm1", "memory", "1024"),
				),
			},
		},
	})
}

func TestVirtualMachine_Clone(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "proxmox_virtual_machine" "vm1" {
  node   = "pve"
  id     = 888
  cores  = 2
  memory = 4096
  clone = 100
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_virtual_machine.vm1", "cores", "2"),
					resource.TestCheckResourceAttr("proxmox_virtual_machine.vm1", "memory", "4096"),
				),
			},
			{
				ResourceName:      "proxmox_virtual_machine.vm1",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "pve,888",
			},
			{
				Config: providerConfig + `
resource "proxmox_virtual_machine" "vm1" {
  node   = "pve"
  id     = 888
  cores  = 1
  memory = 1024
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_virtual_machine.vm1", "cores", "1"),
					resource.TestCheckResourceAttr("proxmox_virtual_machine.vm1", "memory", "1024"),
				),
			},
		},
	})
}
