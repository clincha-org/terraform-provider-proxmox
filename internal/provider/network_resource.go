package provider

import (
	"context"
	"fmt"

	"github.com/clincha-org/proxmox-api/pkg/proxmox"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var (
	_ resource.Resource              = &networkResource{}
	_ resource.ResourceWithConfigure = &networkResource{}
)

type NetworkResourceModel struct {
	Gateway     string   `tfsdk:"gateway"`
	Type        string   `tfsdk:"type"`
	Autostart   int      `tfsdk:"autostart"`
	Families    []string `tfsdk:"families"`
	Method6     string   `tfsdk:"method6"`
	Iface       string   `tfsdk:"iface"`
	BridgeFd    string   `tfsdk:"bridge_fd"`
	Netmask     string   `tfsdk:"netmask"`
	Priority    int      `tfsdk:"priority"`
	Active      int      `tfsdk:"active"`
	Method      string   `tfsdk:"method"`
	BridgeStp   string   `tfsdk:"bridge_stp"`
	Address     string   `tfsdk:"address"`
	Cidr        string   `tfsdk:"cidr"`
	BridgePorts string   `tfsdk:"bridge_ports"`
}

type networkResource struct {
	client *proxmox.Client
}

func (r networkResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}

	client, ok := request.ProviderData.(*proxmox.Client)

	if !ok {
		response.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *proxmox.Client, got %T. Please report this issue to the developers", request.ProviderData),
		)
	}

	r.client = client
}

func NewNetworkResource() resource.Resource {
	return &networkResource{}
}

func (r networkResource) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_network"
}

func (r networkResource) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"Gateway": schema.StringAttribute{
				Computed: true,
			},
			"Type": schema.StringAttribute{
				Computed: true,
			},
			"Autostart": schema.Int64Attribute{
				Computed: true,
			},
			"Families": schema.ListAttribute{
				Computed: true,
			},
			"Method6": schema.StringAttribute{
				Computed: true,
			},
			"Iface": schema.StringAttribute{
				Computed: true,
			},
			"BridgeFd": schema.StringAttribute{
				Computed: true,
			},
			"Netmask": schema.StringAttribute{
				Computed: true,
			},
			"Priority": schema.Int64Attribute{
				Computed: true,
			},
			"Active": schema.Int64Attribute{
				Computed: true,
			},
			"Method": schema.StringAttribute{
				Computed: true,
			},
			"BridgeStp": schema.StringAttribute{
				Computed: true,
			},
			"Address": schema.StringAttribute{
				Computed: true,
			},
			"Cidr": schema.StringAttribute{
				Computed: true,
			},
			"BridgePorts": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r networkResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan NetworkResourceModel
	diags := request.Plan.Get(ctx, &plan)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	var network = proxmox.Network{
		Gateway:     plan.Gateway,
		Type:        plan.Type,
		Autostart:   plan.Autostart,
		Families:    plan.Families,
		Method6:     plan.Method6,
		Iface:       plan.Iface,
		BridgeFd:    plan.BridgeFd,
		Netmask:     plan.Netmask,
		Priority:    plan.Priority,
		Active:      plan.Active,
		Method:      plan.Method,
		BridgeStp:   plan.BridgeStp,
		Address:     plan.Address,
		Cidr:        plan.Cidr,
		BridgePorts: plan.BridgePorts,
	}

	net, err := r.client.CreateNetwork()

}

func (r networkResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	//TODO implement me
	panic("implement me")
}

func (r networkResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	//TODO implement me
	panic("implement me")
}

func (r networkResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	//TODO implement me
	panic("implement me")
}
