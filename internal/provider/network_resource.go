package provider

import (
	"context"
	"fmt"
	"github.com/clincha-org/proxmox-api/pkg/proxmox"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &networkResource{}
	_ resource.ResourceWithConfigure = &networkResource{}
)

type NetworkResourceModel struct {
	Interface              types.String `tfsdk:"iface"`
	Type                   types.String `tfsdk:"type"`
	Address                types.String `tfsdk:"address"`
	Address6               types.String `tfsdk:"address6"`
	Autostart              types.Int64  `tfsdk:"autostart"`
	BondPrimary            types.String `tfsdk:"bond_primary"`
	BondMode               types.String `tfsdk:"bond_mode"`
	BondMiiMon             types.String `tfsdk:"bond_miimon"`
	BondTransmitHashPolicy types.String `tfsdk:"bond_xmit_hash_policy"`
	BridgePorts            types.String `tfsdk:"bridge_ports"`
	BridgeVlanAware        types.Int64  `tfsdk:"bridge_vlan_aware"`
	CIDR                   types.String `tfsdk:"cidr"`
	CIDR6                  types.String `tfsdk:"cidr6"`
	Comments               types.String `tfsdk:"comments"`
	Comments6              types.String `tfsdk:"comments6"`
	Gateway                types.String `tfsdk:"gateway"`
	Gateway6               types.String `tfsdk:"gateway6"`
	MTU                    types.Int64  `tfsdk:"mtu"`
	Netmask                types.String `tfsdk:"netmask"`
	Netmask6               types.String `tfsdk:"netmask6"`
	OVSBonds               types.String `tfsdk:"ovs_bonds"`
	OVSBridge              types.String `tfsdk:"ovs_bridge"`
	OVSOptions             types.String `tfsdk:"ovs_options"`
	OVSPorts               types.String `tfsdk:"ovs_ports"`
	OVSTag                 types.String `tfsdk:"ovs_tag"`
	Slaves                 types.String `tfsdk:"slaves"`
	VlanID                 types.Int64  `tfsdk:"vlan_id"`
	VlanRawDevice          types.String `tfsdk:"vlan_raw_device"`
	Families               types.List   `tfsdk:"families"`
	Method                 types.String `tfsdk:"method"`
	Method6                types.String `tfsdk:"method6"`
	BridgeFd               types.String `tfsdk:"bridge_fd"`
	Priority               types.Int64  `tfsdk:"priority"`
	Active                 types.Int64  `tfsdk:"active"`
	BridgeStp              types.String `tfsdk:"bridge_stp"`
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
			"iface": schema.StringAttribute{
				Required: true,
			},
			"type": schema.StringAttribute{
				Required: true,
			},
			"address": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"address6": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"autostart": schema.Int64Attribute{
				Optional: true,
				Computed: true,
			},
			"bond_primary": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"bond_mode": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"bond_miimon": schema.StringAttribute{
				Computed: true,
			},
			"bond_xmit_hash_policy": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"bridge_ports": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"bridge_vlan_aware": schema.Int64Attribute{
				Optional: true,
				Computed: true,
			},
			"cidr": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"cidr6": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"comments": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"comments6": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"gateway": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"gateway6": schema.StringAttribute{
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
			"netmask6": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"ovs_bonds": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"ovs_bridge": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"ovs_options": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"ovs_ports": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"ovs_tag": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"slaves": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"vlan_id": schema.Int64Attribute{
				Optional: true,
				Computed: true,
			},
			"vlan_raw_device": schema.StringAttribute{
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
			"method6": schema.StringAttribute{
				Computed: true,
			},
			"bridge_fd": schema.StringAttribute{
				Computed: true,
			},
			"priority": schema.Int64Attribute{
				Computed: true,
			},
			"active": schema.Int64Attribute{
				Computed: true,
			},
			"bridge_stp": schema.StringAttribute{
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
		Interface:              plan.Interface.ValueString(),
		Type:                   plan.Type.ValueString(),
		Address:                plan.Address.ValueString(),
		Address6:               plan.Address6.ValueString(),
		AutoStart:              plan.Autostart.ValueInt64(),
		BondPrimary:            plan.BondPrimary.ValueString(),
		BondMode:               plan.BondMode.ValueString(),
		BondTransmitHashPolicy: plan.BondTransmitHashPolicy.ValueString(),
		BridgePorts:            plan.BridgePorts.ValueString(),
		BridgeVlanAware:        plan.BridgeVlanAware.ValueInt64(),
		CIDR:                   plan.CIDR.ValueString(),
		CIDR6:                  plan.CIDR6.ValueString(),
		Comments:               plan.Comments.ValueString(),
		Comments6:              plan.Comments6.ValueString(),
		Gateway:                plan.Gateway.ValueString(),
		Gateway6:               plan.Gateway6.ValueString(),
		MTU:                    plan.MTU.ValueInt64(),
		Netmask:                plan.Netmask.ValueString(),
		Netmask6:               plan.Netmask6.ValueString(),
		OVSBonds:               plan.OVSBonds.ValueString(),
		OVSBridge:              plan.OVSBridge.ValueString(),
		OVSOptions:             plan.OVSOptions.ValueString(),
		OVSPorts:               plan.OVSPorts.ValueString(),
		OVSTag:                 plan.OVSTag.ValueString(),
		Slaves:                 plan.Slaves.ValueString(),
		VlanID:                 plan.VlanID.ValueInt64(),
		VlanRawDevice:          plan.VlanRawDevice.ValueString(),
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
		Interface:              types.StringValue(network.Interface),
		Type:                   types.StringValue(network.Type),
		Address:                types.StringValue(network.Address),
		Address6:               types.StringValue(network.Address6),
		Autostart:              types.Int64Value(network.Autostart),
		BondPrimary:            types.StringValue(network.BondPrimary),
		BondMode:               types.StringValue(network.BondMode),
		BondMiiMon:             types.StringValue(network.BondMiiMon),
		BondTransmitHashPolicy: types.StringValue(network.BondTransmitHashPolicy),
		BridgePorts:            types.StringValue(network.BridgePorts),
		BridgeVlanAware:        types.Int64Value(network.BridgeVlanAware),
		CIDR:                   types.StringValue(network.CIDR),
		CIDR6:                  types.StringValue(network.CIDR6),
		Comments:               types.StringValue(network.Comments),
		Comments6:              types.StringValue(network.Comments6),
		Gateway:                types.StringValue(network.Gateway),
		Gateway6:               types.StringValue(network.Gateway6),
		MTU:                    types.Int64Value(network.MTU),
		Netmask:                types.StringValue(network.Netmask),
		Netmask6:               types.StringValue(network.Netmask6),
		OVSBonds:               types.StringValue(network.OVSBonds),
		OVSBridge:              types.StringValue(network.OVSBridge),
		OVSOptions:             types.StringValue(network.OVSOptions),
		OVSPorts:               types.StringValue(network.OVSPorts),
		OVSTag:                 types.StringValue(network.OVSTag),
		Slaves:                 types.StringValue(network.Slaves),
		VlanID:                 types.Int64Value(network.VlanID),
		VlanRawDevice:          types.StringValue(network.VlanRawDevice),
		Method:                 types.StringValue(network.Method),
		Method6:                types.StringValue(network.Method6),
		BridgeFd:               types.StringValue(network.BridgeFd),
		Priority:               types.Int64Value(network.Priority),
		Active:                 types.Int64Value(network.Active),
		BridgeStp:              types.StringValue(network.BridgeStp),
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
		Interface:              types.StringValue(network.Interface),
		Type:                   types.StringValue(network.Type),
		Address:                types.StringValue(network.Address),
		Address6:               types.StringValue(network.Address6),
		Autostart:              types.Int64Value(network.Autostart),
		BondPrimary:            types.StringValue(network.BondPrimary),
		BondMode:               types.StringValue(network.BondMode),
		BondMiiMon:             types.StringValue(network.BondMiiMon),
		BondTransmitHashPolicy: types.StringValue(network.BondTransmitHashPolicy),
		BridgePorts:            types.StringValue(network.BridgePorts),
		BridgeVlanAware:        types.Int64Value(network.BridgeVlanAware),
		CIDR:                   types.StringValue(network.CIDR),
		CIDR6:                  types.StringValue(network.CIDR6),
		Comments:               types.StringValue(network.Comments),
		Comments6:              types.StringValue(network.Comments6),
		Gateway:                types.StringValue(network.Gateway),
		Gateway6:               types.StringValue(network.Gateway6),
		MTU:                    types.Int64Value(network.MTU),
		Netmask:                types.StringValue(network.Netmask),
		Netmask6:               types.StringValue(network.Netmask6),
		OVSBonds:               types.StringValue(network.OVSBonds),
		OVSBridge:              types.StringValue(network.OVSBridge),
		OVSOptions:             types.StringValue(network.OVSOptions),
		OVSPorts:               types.StringValue(network.OVSPorts),
		OVSTag:                 types.StringValue(network.OVSTag),
		Slaves:                 types.StringValue(network.Slaves),
		VlanID:                 types.Int64Value(network.VlanID),
		VlanRawDevice:          types.StringValue(network.VlanRawDevice),
		Method:                 types.StringValue(network.Method),
		Method6:                types.StringValue(network.Method6),
		BridgeFd:               types.StringValue(network.BridgeFd),
		Priority:               types.Int64Value(network.Priority),
		Active:                 types.Int64Value(network.Active),
		BridgeStp:              types.StringValue(network.BridgeStp),
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

}

func (r *networkResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	//TODO implement me
	panic("implement me")
}
