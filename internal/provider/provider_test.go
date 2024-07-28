package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	// providerConfig is a shared configuration to combine with the actual
	// test configuration so the Proxmox client is properly configured.
	// It is also possible to use the PROXMOX_ environment variables instead,
	// such as updating the Makefile and running the testing through that tool.
	providerConfig = `
provider "proxmox" {
  username = "root@pam"
  password = "vagrant"
  host     = "https://localhost:8006"
}
`
)

var (
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"proxmox": providerserver.NewProtocol6WithError(New("test")()),
	}
)
