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
resource "proxmox_network_bridge" "vmbr88" {
  interface = "vmbr88"
  address   = "192.168.1.88"
  netmask   = "255.255.255.0"
  autostart = true
  bridge_vlan_aware = true
  comments  = "Test network"
  mtu = 1500
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_network_bridge.vmbr88", "interface", "vmbr88"),
					resource.TestCheckResourceAttr("proxmox_network_bridge.vmbr88", "address", "192.168.1.88"),
					resource.TestCheckResourceAttr("proxmox_network_bridge.vmbr88", "netmask", "255.255.255.0"),
					resource.TestCheckResourceAttr("proxmox_network_bridge.vmbr88", "autostart", "true"),
					resource.TestCheckResourceAttr("proxmox_network_bridge.vmbr88", "comments", "Test network"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "proxmox_network_bridge.vmbr88",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "proxmox_network_bridge" "vmbr88" {
  interface = "vmbr88"
  address   = "192.168.1.88"
  netmask   = "255.255.255.0"
  bridge_vlan_aware = true
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_network_bridge.vmbr88", "interface", "vmbr88"),
					resource.TestCheckResourceAttr("proxmox_network_bridge.vmbr88", "mtu", "1500"),
					resource.TestCheckResourceAttr("proxmox_network_bridge.vmbr88", "autostart", "true"),
					resource.TestCheckResourceAttr("proxmox_network_bridge.vmbr88", "comments", "Test network"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
