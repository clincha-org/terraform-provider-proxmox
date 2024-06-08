package provider

import (
	"context"
	"fmt"
	"github.com/clincha-org/proxmox-api/pkg/proxmox"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource              = &networkResource{}
	_ resource.ResourceWithConfigure = &networkResource{}
)

type NetworkResourceModel struct {
	Interface       types.String `tfsdk:"interface"`
	Type            types.String `tfsdk:"type"`
	Address         types.String `tfsdk:"address"`
	Autostart       types.Bool   `tfsdk:"autostart"`
	BridgePorts     types.String `tfsdk:"bridge_ports"`
	BridgeVlanAware types.Bool   `tfsdk:"bridge_vlan_aware"`
	CIDR            types.String `tfsdk:"cidr"`
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

func (r *networkResource) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_network"
}

func (r *networkResource) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
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
			"cidr": schema.StringAttribute{
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
			"vlan_id": schema.Int64Attribute{
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

func (r *networkResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan NetworkResourceModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	autostart := plan.Autostart.ValueBool()
	bridgeVlanAware := plan.BridgeVlanAware.ValueBool()
	networkRequest := proxmox.NetworkRequest{
		Interface:       plan.Interface.ValueString(),
		Type:            plan.Type.ValueString(),
		Address:         plan.Address.ValueString(),
		AutoStart:       &autostart,
		BridgePorts:     plan.BridgePorts.ValueString(),
		BridgeVlanAware: &bridgeVlanAware,
		CIDR:            plan.CIDR.ValueString(),
		Comments:        plan.Comments.ValueString(),
		Gateway:         plan.Gateway.ValueString(),
		MTU:             plan.MTU.ValueInt64(),
		Netmask:         plan.Netmask.ValueString(),
		VlanID:          plan.VlanID.ValueInt64(),
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
		Interface:       types.StringValue(network.Interface),
		Type:            types.StringValue(network.Type),
		Address:         types.StringValue(network.Address),
		Autostart:       types.BoolValue(network.Autostart == 1),
		BridgePorts:     types.StringValue(network.BridgePorts),
		BridgeVlanAware: types.BoolValue(network.BridgeVlanAware == 1),
		CIDR:            types.StringValue(network.CIDR),
		Comments:        types.StringValue(network.Comments),
		Gateway:         types.StringValue(network.Gateway),
		MTU:             types.Int64Value(network.MTU),
		Netmask:         types.StringValue(network.Netmask),
		VlanID:          types.Int64Value(network.VlanID),
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
		Interface:       types.StringValue(network.Interface),
		Type:            types.StringValue(network.Type),
		Address:         types.StringValue(network.Address),
		Autostart:       types.BoolValue(network.Autostart == 1),
		BridgePorts:     types.StringValue(network.BridgePorts),
		BridgeVlanAware: types.BoolValue(network.BridgeVlanAware == 1),
		CIDR:            types.StringValue(network.CIDR),
		Comments:        types.StringValue(network.Comments),
		Gateway:         types.StringValue(network.Gateway),
		MTU:             types.Int64Value(network.MTU),
		Netmask:         types.StringValue(network.Netmask),
		VlanID:          types.Int64Value(network.VlanID),
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

	autostart := plan.Autostart.ValueBool()
	bridgeVlanAware := plan.BridgeVlanAware.ValueBool()
	networkRequest := proxmox.NetworkRequest{
		Interface:       plan.Interface.ValueString(),
		Type:            plan.Type.ValueString(),
		Address:         plan.Address.ValueString(),
		AutoStart:       &autostart,
		BridgePorts:     plan.BridgePorts.ValueString(),
		BridgeVlanAware: &bridgeVlanAware,
		CIDR:            plan.CIDR.ValueString(),
		Comments:        plan.Comments.ValueString(),
		Gateway:         plan.Gateway.ValueString(),
		MTU:             plan.MTU.ValueInt64(),
		Netmask:         plan.Netmask.ValueString(),
		VlanID:          plan.VlanID.ValueInt64(),
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
		Interface:       types.StringValue(network.Interface),
		Type:            types.StringValue(network.Type),
		Address:         types.StringValue(network.Address),
		Autostart:       types.BoolValue(network.Autostart == 1),
		BridgePorts:     types.StringValue(network.BridgePorts),
		BridgeVlanAware: types.BoolValue(network.BridgeVlanAware == 1),
		CIDR:            types.StringValue(network.CIDR),
		Comments:        types.StringValue(network.Comments),
		Gateway:         types.StringValue(network.Gateway),
		MTU:             types.Int64Value(network.MTU),
		Netmask:         types.StringValue(network.Netmask),
		VlanID:          types.Int64Value(network.VlanID),
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
