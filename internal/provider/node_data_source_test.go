package provider

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestNodeDataSource_Read(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "proxmox_node" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.proxmox_node.test", "nodes.0.id", "node/pve"),
					resource.TestCheckResourceAttr("data.proxmox_node.test", "nodes.0.type", "node"),
					resource.TestCheckResourceAttr("data.proxmox_node.test", "nodes.0.node", "pve"),
					resource.TestCheckResourceAttr("data.proxmox_node.test", "nodes.0.status", "online"),
				),
			},
		},
	})
}
