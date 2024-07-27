package provider

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestNetworkResource_Create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "proxmox_network" "vmbr88" {
  interface = "vmbr88"
  type      = "bridge"
  address   = "10.1.2.3"
  netmask = "255.255.255.0"
  autostart = false
  comments = "Test network"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_network.vmbr88", "interface", "vmbr88"),
					resource.TestCheckResourceAttr("proxmox_network.vmbr88", "type", "bridge"),
					resource.TestCheckResourceAttr("proxmox_network.vmbr88", "address", "10.1.2.3"),
					resource.TestCheckResourceAttr("proxmox_network.vmbr88", "netmask", "255.255.255.0"),
					resource.TestCheckResourceAttr("proxmox_network.vmbr88", "autostart", "false"),
					resource.TestCheckResourceAttr("proxmox_network.vmbr88", "comments", "Test network"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "proxmox_network.vmbr88",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "proxmox_network" "vmbr88" {
  interface = "vmbr88"
  type      = "bridge"
  address   = "0.0.0.0"
  netmask   = "255.255.255.0"
  autostart = false
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_network.vmbr88", "interface", "vmbr88"),
					resource.TestCheckResourceAttr("proxmox_network.vmbr88", "type", "bridge"),
					resource.TestCheckResourceAttr("proxmox_network.vmbr88", "address", ""),
					resource.TestCheckResourceAttr("proxmox_network.vmbr88", "netmask", ""),
					resource.TestCheckResourceAttr("proxmox_network.vmbr88", "autostart", "false"),
					resource.TestCheckResourceAttr("proxmox_network.vmbr88", "comments", ""),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
