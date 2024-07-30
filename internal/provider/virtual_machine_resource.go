package provider

import (
	"context"
	"fmt"
	"github.com/clincha-org/proxmox-api/pkg/proxmox"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &virtualMachineResource{}
	_ resource.ResourceWithConfigure   = &virtualMachineResource{}
	_ resource.ResourceWithImportState = &virtualMachineResource{}
)

type virtualMachineModel struct {
	Node   types.String `tfsdk:"node"`
	ID     types.Int64  `tfsdk:"id"`
	Memory types.String `tfsdk:"memory"`
	Cores  types.Int64  `tfsdk:"cores"`
}

type virtualMachineResource struct {
	client *proxmox.Client
}

func NewVirtualMachineResource() resource.Resource {
	return &virtualMachineResource{}
}

func (v *virtualMachineResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	//TODO implement me
	panic("implement me")
}

func (v *virtualMachineResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}

	client, ok := request.ProviderData.(*proxmox.Client)
	if !ok {
		response.Diagnostics.AddError(
			"virtual_machine_configure",
			fmt.Sprintf("Expected *proxmox.Client, got %T", request.ProviderData),
		)
		return
	}

	v.client = client
}

func (v *virtualMachineResource) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_virtual_machine"
}

func (v *virtualMachineResource) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"node": schema.StringAttribute{
				Required:    true,
				Description: "The node where the virtual machine is located",
			},
			"id": schema.Int64Attribute{
				Required:    true,
				Description: "The virtual machine ID",
			},
			"memory": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The amount of memory for the virtual machine",
			},
			"cores": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The number of cores for the virtual machine",
			},
		},
	}
}

func (v *virtualMachineResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var data virtualMachineModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	vmRequest := proxmox.VirtualMachineRequest{
		ID:           data.ID.ValueInt64(),
		Cdrom:        "local:iso/ubuntu-22.04.4-live-server-amd64.iso",
		SCSI1:        "local-lvm:8",
		Net1:         "model=virtio,bridge=vmbr0,firewall=1",
		SCSIHardware: "virtio-scsi-pci",
	}

	vm, err := v.client.CreateVM(data.Node.ValueString(), &vmRequest, true)
	if err != nil {
		return
	}

	state := virtualMachineModel{
		Node:   data.Node,
		ID:     data.ID,
		Memory: types.StringValue(vm.Memory),
		Cores:  types.Int64Value(vm.Cores),
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (v *virtualMachineResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	//TODO implement me
	panic("implement me")
}

func (v *virtualMachineResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	//TODO implement me
	panic("implement me")
}

func (v *virtualMachineResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	//TODO implement me
	panic("implement me")
}
