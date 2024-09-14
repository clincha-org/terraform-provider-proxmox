package provider

import (
	"context"
	"fmt"
	"github.com/clincha-org/proxmox-api/pkg/ide"
	"github.com/clincha-org/proxmox-api/pkg/proxmox"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strconv"
	"strings"
)

var (
	_ resource.Resource                = &virtualMachineResource{}
	_ resource.ResourceWithConfigure   = &virtualMachineResource{}
	_ resource.ResourceWithImportState = &virtualMachineResource{}
)

type virtualMachineModel struct {
	Node       types.String               `tfsdk:"node"`
	ID         types.Int64                `tfsdk:"id"`
	Memory     types.Int64                `tfsdk:"memory"`
	Cores      types.Int64                `tfsdk:"cores"`
	IDEDevices []InternalDataStorageModel `tfsdk:"ide_devices"`
}

type InternalDataStorageModel struct {
	ID      types.Int64  `tfsdk:"id"`
	Storage types.String `tfsdk:"storage"`
	Path    types.String `tfsdk:"path"`
	Media   types.String `tfsdk:"media"`
	Size    types.String `tfsdk:"size"`
}

type virtualMachineResource struct {
	client *proxmox.Client
}

func NewVirtualMachineResource() resource.Resource {
	return &virtualMachineResource{}
}

func (v *virtualMachineResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	// The ID is a combination of the name of the node and the id of the virtual machine
	idParts := strings.Split(request.ID, ",")
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		response.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Please provive the identifier in the format: node,vmid. For example: pve,123. Got: %q", request.ID),
		)
		return
	}
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("node"), idParts[0])...)

	// The ID needs to be converted into an integer
	identifier, err := strconv.ParseInt(idParts[1], 10, 64)
	if err != nil {
		response.Diagnostics.AddError(
			"Failed to convert import identifier",
			fmt.Sprintf("Unable to convert the ID component of the import identifier to an integer: %s", idParts[1]),
		)
		return
	}
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("id"), identifier)...)
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
			"memory": schema.Int64Attribute{
				Required:    true,
				Description: "The amount of memory for the virtual machine",
			},
			"cores": schema.Int64Attribute{
				Required:    true,
				Description: "The number of cores for the virtual machine",
			},
			"ide_devices": schema.ListNestedAttribute{
				Optional:    true,
				Description: "The IDE devices for the virtual machine",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Required:    true,
							Description: "The ID of the IDE device",
						},
						"storage": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "The storage for the IDE device",
						},
						"path": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "The path for the IDE device",
						},
						"media": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "The media for the IDE device",
						},
						"size": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "The size for the IDE device",
						},
					},
				},
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

	vm := proxmox.VirtualMachine{
		ID:     data.ID.ValueInt64(),
		Cores:  data.Cores.ValueInt64(),
		Memory: data.Memory.ValueInt64(),
	}

	if data.IDEDevices != nil {
		var devices []ide.InternalDataStorage
		for _, ideDevice := range data.IDEDevices {
			device := ide.InternalDataStorage{
				ID:      ideDevice.ID.ValueInt64(),
				Storage: ideDevice.Storage.ValueString(),
				Path:    ideDevice.Path.ValueStringPointer(),
				Media:   ideDevice.Media.ValueStringPointer(),
				Size:    ideDevice.Size.ValueStringPointer(),
			}

			if ideDevice.Path.IsUnknown() {
				device.Path = nil
			}
			if ideDevice.Media.IsUnknown() {
				device.Media = nil
			}
			if ideDevice.Size.IsUnknown() {
				device.Size = nil
			}

			devices = append(devices, device)
		}
		vm.IDEDevices = &devices
	}

	vm, err := v.client.CreateVM(data.Node.ValueString(), &vm, true)
	if err != nil {
		response.Diagnostics.AddError("virtual_machine_create", err.Error())
		return
	}

	state := virtualMachineModel{
		Node:   data.Node,
		ID:     data.ID,
		Memory: types.Int64Value(vm.Memory),
		Cores:  types.Int64Value(vm.Cores),
	}

	if vm.IDEDevices != nil {
		var tfdevices []InternalDataStorageModel
		for _, device := range *vm.IDEDevices {
			tfdevice := InternalDataStorageModel{
				ID:      types.Int64Value(device.ID),
				Storage: types.StringValue(device.Storage),
				Path:    types.StringPointerValue(device.Path),
				Media:   types.StringPointerValue(device.Media),
				Size:    types.StringPointerValue(device.Size),
			}
			tfdevices = append(tfdevices, tfdevice)
		}
		state.IDEDevices = tfdevices
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (v *virtualMachineResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var expectedState virtualMachineModel
	response.Diagnostics.Append(request.State.Get(ctx, &expectedState)...)
	if response.Diagnostics.HasError() {
		return
	}

	vm, err := v.client.GetVM(expectedState.Node.ValueString(), expectedState.ID.ValueInt64())
	if err != nil {
		response.Diagnostics.AddError("virtual_machine_read", err.Error())
		return
	}

	state := virtualMachineModel{
		Node:   expectedState.Node,
		ID:     expectedState.ID,
		Memory: types.Int64Value(vm.Memory),
		Cores:  types.Int64Value(vm.Cores),
	}

	if vm.IDEDevices != nil {
		var tfdevices []InternalDataStorageModel
		for _, device := range *vm.IDEDevices {
			tfdevice := InternalDataStorageModel{
				ID:      types.Int64Value(device.ID),
				Storage: types.StringValue(device.Storage),
				Path:    types.StringPointerValue(device.Path),
				Media:   types.StringPointerValue(device.Media),
				Size:    types.StringPointerValue(device.Size),
			}
			tfdevices = append(tfdevices, tfdevice)
		}
		state.IDEDevices = tfdevices
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (v *virtualMachineResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var data virtualMachineModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)
	if response.Diagnostics.HasError() {
		return
	}

	vm := proxmox.VirtualMachine{
		ID:     data.ID.ValueInt64(),
		Cores:  data.Cores.ValueInt64(),
		Memory: data.Memory.ValueInt64(),
	}

	if data.IDEDevices != nil {
		var devices []ide.InternalDataStorage
		for _, ideDevice := range data.IDEDevices {
			device := ide.InternalDataStorage{
				ID:      ideDevice.ID.ValueInt64(),
				Storage: ideDevice.Storage.ValueString(),
				Path:    ideDevice.Path.ValueStringPointer(),
				Media:   ideDevice.Media.ValueStringPointer(),
				Size:    ideDevice.Size.ValueStringPointer(),
			}

			if ideDevice.Path.IsUnknown() {
				device.Path = nil
			}
			if ideDevice.Media.IsUnknown() {
				device.Media = nil
			}
			if ideDevice.Size.IsUnknown() {
				device.Size = nil
			}

			devices = append(devices, device)
		}
		vm.IDEDevices = &devices
	}

	vm, err := v.client.UpdateVM(data.Node.ValueString(), &vm)
	if err != nil {
		response.Diagnostics.AddError("virtual_machine_update", err.Error())
		return
	}

	state := virtualMachineModel{
		Node:   data.Node,
		ID:     data.ID,
		Memory: types.Int64Value(vm.Memory),
		Cores:  types.Int64Value(vm.Cores),
	}

	if vm.IDEDevices != nil {
		var tfdevices []InternalDataStorageModel
		for _, device := range *vm.IDEDevices {
			tfdevice := InternalDataStorageModel{
				ID:      types.Int64Value(device.ID),
				Storage: types.StringValue(device.Storage),
				Path:    types.StringPointerValue(device.Path),
				Media:   types.StringPointerValue(device.Media),
				Size:    types.StringPointerValue(device.Size),
			}
			tfdevices = append(tfdevices, tfdevice)
		}
		state.IDEDevices = tfdevices
	}

	response.Diagnostics.Append(response.State.Set(ctx, &state)...)
}

func (v *virtualMachineResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state virtualMachineModel
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	err := v.client.DeleteVM(state.Node.ValueString(), state.ID.ValueInt64())
	if err != nil {
		response.Diagnostics.AddError("virtual_machine_delete", err.Error())
	}
}
