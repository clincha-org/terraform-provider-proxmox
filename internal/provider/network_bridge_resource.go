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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

var (
	_ resource.Resource                = &networkBridgeResource{}
	_ resource.ResourceWithConfigure   = &networkBridgeResource{}
	_ resource.ResourceWithImportState = &networkBridgeResource{}
)

type NetworkBridgeResourceModel struct {
	ID              types.String `tfsdk:"id"`
	Node            types.String `tfsdk:"node"`
	Interface       types.String `tfsdk:"interface"`
	Address         types.String `tfsdk:"address"`
	Autostart       types.Bool   `tfsdk:"autostart"`
	BridgePorts     types.String `tfsdk:"bridge_ports"`
	BridgeVlanAware types.Bool   `tfsdk:"bridge_vlan_aware"`
	Comments        types.String `tfsdk:"comments"`
	Gateway         types.String `tfsdk:"gateway"`
	MTU             types.Int64  `tfsdk:"mtu"`
	Netmask         types.String `tfsdk:"netmask"`
	Families        types.List   `tfsdk:"families"`
	Method          types.String `tfsdk:"method"`
	Active          types.Bool   `tfsdk:"active"`
}

type networkBridgeResource struct {
	client *proxmox.Client
}

func NewNetworkResource() resource.Resource {
	return &networkBridgeResource{}
}

func (r *networkBridgeResource) Metadata(_ context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_network_bridge"
}

func (r *networkBridgeResource) Configure(_ context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (r *networkBridgeResource) Schema(_ context.Context, _ resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"node": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"interface": schema.StringAttribute{
				Required: true,
			},
			"address": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"autostart": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"bridge_ports": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"bridge_vlan_aware": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"comments": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"gateway": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"mtu": schema.Int64Attribute{
				Optional: true,
				Computed: true,
			},
			"netmask": schema.StringAttribute{
				Optional: true,
				Computed: true,
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

func (r *networkBridgeResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	// The ID is a combination of the name of the node and the name of the interface
	idParts := strings.Split(request.ID, ",")
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		response.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Please provive the identifier in the format: node,interface. For example: pve,vmbr0. Got: %q", request.ID),
		)
		return
	}
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("node"), idParts[0])...)
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("interface"), idParts[1])...)
}

func (r *networkBridgeResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan NetworkBridgeResourceModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	networkRequest := proxmox.NetworkRequest{
		Interface:       plan.Interface.ValueString(),
		Type:            "bridge",
		Address:         plan.Address.ValueStringPointer(),
		AutoStart:       plan.Autostart.ValueBoolPointer(),
		BridgePorts:     plan.BridgePorts.ValueStringPointer(),
		BridgeVlanAware: plan.BridgeVlanAware.ValueBoolPointer(),
		Comments:        plan.Comments.ValueStringPointer(),
		Gateway:         plan.Gateway.ValueStringPointer(),
		MTU:             plan.MTU.ValueInt64Pointer(),
		Netmask:         plan.Netmask.ValueStringPointer(),
	}

	// We need to set the fields to nil if they are unknown
	// This is because the Proxmox API will interpret an empty string as a value
	// that will cause a value to be set when the intention is to omit the value.
	if plan.Address.IsUnknown() {
		networkRequest.Address = nil
	}
	if plan.BridgePorts.IsUnknown() {
		networkRequest.BridgePorts = nil
	}
	if plan.BridgeVlanAware.IsUnknown() {
		networkRequest.BridgeVlanAware = nil
	}
	if plan.Comments.IsUnknown() {
		networkRequest.Comments = nil
	}
	if plan.Gateway.IsUnknown() {
		networkRequest.Gateway = nil
	}
	if plan.MTU.IsUnknown() {
		networkRequest.MTU = nil
	}
	if plan.Netmask.IsUnknown() {
		networkRequest.Netmask = nil
	}

	node := proxmox.Node{Node: plan.Node.ValueString()}
	network, err := r.client.CreateNetwork(&node, &networkRequest)
	if err != nil {
		response.Diagnostics.AddError(
			"Error creating Proxmox network",
			"Could not create the Proxmox network: "+plan.Interface.ValueString()+": "+err.Error(),
		)
		return
	}

	state := NetworkBridgeResourceModel{
		ID:              types.StringValue(network.Interface),
		Node:            plan.Node,
		Interface:       types.StringValue(network.Interface),
		Address:         types.StringPointerValue(network.Address),
		Autostart:       types.BoolValue(network.Autostart == 1),
		BridgePorts:     types.StringPointerValue(network.BridgePorts),
		BridgeVlanAware: types.BoolValue(network.BridgeVlanAware == 1),
		Comments:        types.StringPointerValue(network.Comments),
		Gateway:         types.StringPointerValue(network.Gateway),
		MTU:             types.Int64PointerValue(network.MTU),
		Netmask:         types.StringPointerValue(network.Netmask),
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
}

func (r *networkBridgeResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state NetworkBridgeResourceModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	node := proxmox.Node{Node: state.Node.ValueString()}
	network, err := r.client.GetNetwork(&node, state.Interface.ValueString())
	if err != nil {
		response.Diagnostics.AddError(
			"Error reading Proxmox network",
			"Could not read the Proxmox network: "+state.Interface.ValueString()+": "+err.Error(),
		)
		return
	}

	state = NetworkBridgeResourceModel{
		ID:              types.StringValue(network.Interface),
		Node:            types.StringValue(node.Node),
		Interface:       types.StringValue(network.Interface),
		Address:         types.StringPointerValue(network.Address),
		Autostart:       types.BoolValue(network.Autostart == 1),
		BridgePorts:     types.StringPointerValue(network.BridgePorts),
		BridgeVlanAware: types.BoolValue(network.BridgeVlanAware == 1),
		Comments:        types.StringPointerValue(network.Comments),
		Gateway:         types.StringPointerValue(network.Gateway),
		MTU:             types.Int64PointerValue(network.MTU),
		Netmask:         types.StringPointerValue(network.Netmask),
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
}

func (r *networkBridgeResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan NetworkBridgeResourceModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	networkRequest := proxmox.NetworkRequest{
		Interface:       plan.Interface.ValueString(),
		Type:            "bridge",
		Address:         plan.Address.ValueStringPointer(),
		AutoStart:       plan.Autostart.ValueBoolPointer(),
		BridgePorts:     plan.BridgePorts.ValueStringPointer(),
		BridgeVlanAware: plan.BridgeVlanAware.ValueBoolPointer(),
		Comments:        plan.Comments.ValueStringPointer(),
		Gateway:         plan.Gateway.ValueStringPointer(),
		Netmask:         plan.Netmask.ValueStringPointer(),
		MTU:             plan.MTU.ValueInt64Pointer(),
	}

	// We need to set the fields to nil if they are unknown
	// This is because the Proxmox API will interpret an empty string as a value
	// that will cause a value to be set when the intention is to omit the value.
	if plan.Address.IsUnknown() {
		networkRequest.Address = nil
	}
	if plan.BridgePorts.IsUnknown() {
		networkRequest.BridgePorts = nil
	}
	if plan.BridgeVlanAware.IsUnknown() {
		networkRequest.BridgeVlanAware = nil
	}
	if plan.Comments.IsUnknown() {
		networkRequest.Comments = nil
	}
	if plan.Gateway.IsUnknown() {
		networkRequest.Gateway = nil
	}
	if plan.MTU.IsUnknown() {
		networkRequest.MTU = nil
	}
	if plan.Netmask.IsUnknown() {
		networkRequest.Netmask = nil
	}

	node := proxmox.Node{Node: plan.Node.ValueString()}
	network, err := r.client.UpdateNetwork(&node, &networkRequest)
	if err != nil {
		response.Diagnostics.AddError(
			"Error creating Proxmox network",
			"Could not create the Proxmox network: "+plan.Interface.ValueString()+". Got this error: "+err.Error(),
		)
		return
	}

	plan = NetworkBridgeResourceModel{
		ID:              types.StringValue(network.Interface),
		Node:            types.StringValue(node.Node),
		Interface:       types.StringValue(network.Interface),
		Address:         types.StringPointerValue(network.Address),
		Autostart:       types.BoolValue(network.Autostart == 1),
		BridgePorts:     types.StringPointerValue(network.BridgePorts),
		BridgeVlanAware: types.BoolValue(network.BridgeVlanAware == 1),
		Comments:        types.StringPointerValue(network.Comments),
		Gateway:         types.StringPointerValue(network.Gateway),
		MTU:             types.Int64PointerValue(network.MTU),
		Netmask:         types.StringPointerValue(network.Netmask),
		Method:          types.StringValue(network.Method),
		Active:          types.BoolValue(network.Active == 1),
	}
	var familiesStrings []attr.Value
	for _, family := range network.Families {
		familiesStrings = append(familiesStrings, types.StringValue(family))
	}
	plan.Families, diags = types.ListValue(types.StringType, familiesStrings)
	response.Diagnostics.Append(diags...)

	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
}

func (r *networkBridgeResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var state NetworkBridgeResourceModel
	diags := request.State.Get(ctx, &state)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}
	node := proxmox.Node{Node: state.Node.ValueString()}
	err := r.client.DeleteNetwork(&node, state.Interface.ValueString())
	if err != nil {
		response.Diagnostics.AddError(
			"Error deleting Proxmox network",
			"Could not delete the Proxmox network: "+state.Interface.ValueString()+". Got this error: "+err.Error(),
		)
		return
	}
}
