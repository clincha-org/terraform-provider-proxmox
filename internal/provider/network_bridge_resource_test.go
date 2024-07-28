package provider

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestNetworkResource_Create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "proxmox_network_bridge" "vmbr88" {
  interface = "vmbr88"
  node = "pve"
  address   = "192.168.1.88"
  netmask   = "255.255.255.0"
  autostart = true
  bridge_vlan_aware = true
  comments  = "Test network"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_network_bridge.vmbr88", "interface", "vmbr88"),
					resource.TestCheckResourceAttr("proxmox_network_bridge.vmbr88", "node", "pve"),
					resource.TestCheckResourceAttr("proxmox_network_bridge.vmbr88", "address", "192.168.1.88"),
					resource.TestCheckResourceAttr("proxmox_network_bridge.vmbr88", "netmask", "255.255.255.0"),
					resource.TestCheckResourceAttr("proxmox_network_bridge.vmbr88", "autostart", "true"),
					resource.TestCheckResourceAttr("proxmox_network_bridge.vmbr88", "bridge_vlan_aware", "true"),
					resource.TestCheckResourceAttr("proxmox_network_bridge.vmbr88", "comments", "Test network"),
				),
			},
			{
				ResourceName:      "proxmox_network_bridge.vmbr88",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "pve,vmbr88",
			},
			{
				Config: providerConfig + `
resource "proxmox_network_bridge" "vmbr88" {
  interface = "vmbr88"
  node = "pve"
  autostart = false
  address   = "192.168.1.88"
  netmask   = "255.255.255.254"
  bridge_vlan_aware = true
  mtu = 1900
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("proxmox_network_bridge.vmbr88", "interface", "vmbr88"),
					resource.TestCheckResourceAttr("proxmox_network_bridge.vmbr88", "node", "pve"),
					resource.TestCheckResourceAttr("proxmox_network_bridge.vmbr88", "address", "192.168.1.88"),
					resource.TestCheckResourceAttr("proxmox_network_bridge.vmbr88", "netmask", "255.255.255.254"),
					resource.TestCheckResourceAttr("proxmox_network_bridge.vmbr88", "autostart", "false"),
					resource.TestCheckResourceAttr("proxmox_network_bridge.vmbr88", "bridge_vlan_aware", "true"),
					resource.TestCheckResourceAttr("proxmox_network_bridge.vmbr88", "comments", "Test network"),
					resource.TestCheckResourceAttr("proxmox_network_bridge.vmbr88", "mtu", "1900"),
				),
			},
		},
	})
}
