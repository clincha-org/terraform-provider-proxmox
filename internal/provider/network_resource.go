package provider

import (
	"context"
	"fmt"
	"github.com/clincha-org/proxmox-api/pkg/proxmox"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &networkResource{}
	_ resource.ResourceWithConfigure   = &networkResource{}
	_ resource.ResourceWithImportState = &networkResource{}
)

type NetworkResourceModel struct {
	ID              types.String `tfsdk:"id"`
	Interface       types.String `tfsdk:"interface"`
	Type            types.String `tfsdk:"type"`
	Address         types.String `tfsdk:"address"`
	Autostart       types.Bool   `tfsdk:"autostart"`
	BridgePorts     types.String `tfsdk:"bridge_ports"`
	BridgeVlanAware types.Bool   `tfsdk:"bridge_vlan_aware"`
	Comments        types.String `tfsdk:"comments"`
	Gateway         types.String `tfsdk:"gateway"`
	MTU             types.Int64  `tfsdk:"mtu"`
	Netmask         types.String `tfsdk:"netmask"`
	VlanID          types.Int64  `tfsdk:"vlan_id"`
	Families        types.List   `tfsdk:"families"`
	Method          types.String `tfsdk:"method"`
	Active          types.Bool   `tfsdk:"active"`
}

type networkResource struct {
	client *proxmox.Client
}

func (r *networkResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("interface"), request, response)
}

func (r *networkResource) Configure(_ context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}

	client, ok := request.ProviderData.(*proxmox.Client)

	if !ok {
		response.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *proxmox.Client, got %T. Please report this issue to the developers", request.ProviderData),
		)
		return
	}

	r.client = client
}

func NewNetworkResource() resource.Resource {
	return &networkResource{}
}

func (r *networkResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_network"
}

func (r *networkResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"interface": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				Required: true,
			},
			"address": schema.StringAttribute{
				Optional: true,
				Computed: false,
			},
			"autostart": schema.BoolAttribute{
				Optional: true,
				Computed: false,
			},
			"bridge_ports": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"bridge_vlan_aware": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"comments": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"gateway": schema.StringAttribute{
				Optional: true,
				Computed: false,
			},
			"mtu": schema.Int64Attribute{
				Optional: true,
				Computed: false,
			},
			"netmask": schema.StringAttribute{
				Optional: true,
				Computed: false,
			},
			"vlan_id": schema.Int64Attribute{
				Optional: true,
				Computed: false,
			},
			"families": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"method": schema.StringAttribute{
				Computed: true,
			},
			"active": schema.BoolAttribute{
				Computed: true,
			},
		},
	}
}

func (r *networkResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan NetworkResourceModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	networkRequest := proxmox.NetworkRequest{
		Interface:       plan.Interface.ValueString(),
		Type:            plan.Type.ValueString(),
		Address:         plan.Address.ValueStringPointer(),
		AutoStart:       plan.Autostart.ValueBoolPointer(),
		BridgePorts:     plan.BridgePorts.ValueStringPointer(),
		BridgeVlanAware: plan.BridgeVlanAware.ValueBoolPointer(),
		Comments:        plan.Comments.ValueString(),
		Gateway:         plan.Gateway.ValueStringPointer(),
		MTU:             plan.MTU.ValueInt64Pointer(),
		Netmask:         plan.Netmask.ValueStringPointer(),
		VlanID:          plan.VlanID.ValueInt64Pointer(),
	}

	node := proxmox.Node{
		Node: "pve",
	}

	network, err := r.client.CreateNetwork(&node, &networkRequest)
	if err != nil {
		response.Diagnostics.AddError(
			"Error creating network",
			"Could not create network, unexpected error: "+err.Error(),
		)
		return
	}

	state := NetworkResourceModel{
		ID:              types.StringValue(network.Interface),
		Interface:       types.StringValue(network.Interface),
		Type:            types.StringValue(network.Type),
		Address:         types.StringPointerValue(network.Address),
		Autostart:       types.BoolValue(network.Autostart == 1),
		BridgePorts:     types.StringPointerValue(network.BridgePorts),
		BridgeVlanAware: types.BoolValue(network.BridgeVlanAware == 1),
		Comments:        types.StringValue(network.Comments),
		Gateway:         types.StringPointerValue(network.Gateway),
		MTU:             types.Int64PointerValue(network.MTU),
		Netmask:         types.StringPointerValue(network.Netmask),
		VlanID:          types.Int64PointerValue(network.VlanID),
		Method:          types.StringValue(network.Method),
		Active:          types.BoolValue(network.Active == 1),
	}
	var familiesStrings []attr.Value
	for _, family := range network.Families {
		familiesStrings = append(familiesStrings, types.StringValue(family))
	}
	state.Families, diags = types.ListValue(types.StringType, familiesStrings)
	response.Diagnostics.Append(diags...)

	diags = response.State.Set(ctx, state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *networkResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state NetworkResourceModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	network, err := r.client.GetNetwork(&proxmox.Node{Node: "pve"}, state.Interface.ValueString())
	if err != nil {
		response.Diagnostics.AddError(
			"Error reading Proxmox network",
			"Could not read the Proxmox network: "+state.Interface.ValueString()+": "+err.Error(),
		)
		return
	}

	state = NetworkResourceModel{
		ID:              types.StringValue(network.Interface),
		Interface:       types.StringValue(network.Interface),
		Type:            types.StringValue(network.Type),
		Address:         types.StringPointerValue(network.Address),
		Autostart:       types.BoolValue(network.Autostart == 1),
		BridgePorts:     types.StringPointerValue(network.BridgePorts),
		BridgeVlanAware: types.BoolValue(network.BridgeVlanAware == 1),
		Comments:        types.StringValue(network.Comments),
		Gateway:         types.StringPointerValue(network.Gateway),
		MTU:             types.Int64PointerValue(network.MTU),
		Netmask:         types.StringPointerValue(network.Netmask),
		VlanID:          types.Int64PointerValue(network.VlanID),
		Method:          types.StringValue(network.Method),
		Active:          types.BoolValue(network.Active == 1),
	}
	var familiesStrings []attr.Value
	for _, family := range network.Families {
		familiesStrings = append(familiesStrings, types.StringValue(family))
	}
	state.Families, diags = types.ListValue(types.StringType, familiesStrings)
	response.Diagnostics.Append(diags...)

	diags = response.State.Set(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *networkResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan NetworkResourceModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	networkRequest := proxmox.NetworkRequest{
		Interface:       plan.Interface.ValueString(),
		Type:            plan.Type.ValueString(),
		Address:         plan.Address.ValueStringPointer(),
		AutoStart:       plan.Autostart.ValueBoolPointer(),
		BridgePorts:     plan.BridgePorts.ValueStringPointer(),
		BridgeVlanAware: plan.BridgeVlanAware.ValueBoolPointer(),
		Comments:        plan.Comments.ValueString(),
		Gateway:         plan.Gateway.ValueStringPointer(),
		MTU:             plan.MTU.ValueInt64Pointer(),
		Netmask:         plan.Netmask.ValueStringPointer(),
		VlanID:          plan.VlanID.ValueInt64Pointer(),
	}

	node := proxmox.Node{
		Node: "pve",
	}

	tflog.Debug(ctx, fmt.Sprintf("Plan (Autostart) before network update %+v", plan.Autostart))

	network, err := r.client.UpdateNetwork(&node, &networkRequest)
	if err != nil {
		response.Diagnostics.AddError(
			"Error creating network",
			"Could not create network, unexpected error: "+err.Error(),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Network (Autostart) after network update %+v", network.Autostart))

	plan = NetworkResourceModel{
		ID:              types.StringValue(network.Interface),
		Interface:       types.StringValue(network.Interface),
		Type:            types.StringValue(network.Type),
		Address:         types.StringPointerValue(network.Address),
		Autostart:       types.BoolValue(network.Autostart == 1),
		BridgePorts:     types.StringPointerValue(network.BridgePorts),
		BridgeVlanAware: types.BoolValue(network.BridgeVlanAware == 1),
		Comments:        types.StringValue(network.Comments),
		Gateway:         types.StringPointerValue(network.Gateway),
		MTU:             types.Int64PointerValue(network.MTU),
		Netmask:         types.StringPointerValue(network.Netmask),
		VlanID:          types.Int64PointerValue(network.VlanID),
		Method:          types.StringValue(network.Method),
		Active:          types.BoolValue(network.Active == 1),
	}
	var familiesStrings []attr.Value
	for _, family := range network.Families {
		familiesStrings = append(familiesStrings, types.StringValue(family))
	}
	plan.Families, diags = types.ListValue(types.StringType, familiesStrings)
	response.Diagnostics.Append(diags...)

	tflog.Debug(ctx, fmt.Sprintf("State (Autostart) after state set %+v", plan.Autostart))

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *networkResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state NetworkResourceModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteNetwork(&proxmox.Node{Node: "pve"}, state.Interface.ValueString())
	if err != nil {
		response.Diagnostics.AddError(
			"Error deleting Proxmox network",
			"Could not delete the Proxmox network: "+state.Interface.ValueString()+": "+err.Error(),
		)
		return
	}
}
