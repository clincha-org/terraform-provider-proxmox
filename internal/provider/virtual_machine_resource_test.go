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
