package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/clincha-org/proxmox-api/pkg/proxmox"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var (
	_ datasource.DataSource              = &nodeDataSource{}
	_ datasource.DataSourceWithConfigure = &nodeDataSource{}
)

type nodeDataSource struct {
	client *proxmox.Client
}

type NodeModel struct {
	Nodes []Node `tfsdk:"nodes"`
}

type Node struct {
	Type           types.String  `tfsdk:"type"`
	Maxcpu         types.Int64   `tfsdk:"max_cpu"`
	Cpu            types.Float64 `tfsdk:"cpu"`
	Status         types.String  `tfsdk:"status"`
	Maxmem         types.Int64   `tfsdk:"max_memory"`
	SslFingerprint types.String  `tfsdk:"ssl_fingerprint"`
	Mem            types.Int64   `tfsdk:"memory"`
	Id             types.String  `tfsdk:"id"`
	Node           types.String  `tfsdk:"node"`
	Disk           types.Int64   `tfsdk:"disk"`
	Uptime         types.Int64   `tfsdk:"uptime"`
	Maxdisk        types.Int64   `tfsdk:"max_disk"`
	Level          types.String  `tfsdk:"level"`
}

func NewNodeDataSource() datasource.DataSource {
	return &nodeDataSource{}
}

func (d *nodeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_node"
}

func (d *nodeDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"nodes": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Computed: true,
						},
						"max_cpu": schema.Int64Attribute{
							Computed: true,
						},
						"cpu": schema.Float64Attribute{
							Computed: true,
						},
						"status": schema.StringAttribute{
							Computed: true,
						},
						"max_memory": schema.Int64Attribute{
							Computed: true,
						},
						"ssl_fingerprint": schema.StringAttribute{
							Computed: true,
						},
						"memory": schema.Int64Attribute{
							Computed: true,
						},
						"id": schema.StringAttribute{
							Computed: true,
						},
						"node": schema.StringAttribute{
							Computed: true,
						},
						"disk": schema.Int64Attribute{
							Computed: true,
						},
						"uptime": schema.Int64Attribute{
							Computed: true,
						},
						"max_disk": schema.Int64Attribute{
							Computed: true,
						},
						"level": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (d *nodeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state NodeModel

	nodes, err := d.client.GetNodes()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Proxmox nodes",
			err.Error(),
		)
		return
	}

	for _, node := range nodes {
		nodeState := Node{
			Type:           types.StringValue(node.Type),
			Maxcpu:         types.Int64Value(int64(node.Maxcpu)),
			Cpu:            types.Float64Value(node.Cpu),
			Status:         types.StringValue(node.Status),
			Maxmem:         types.Int64Value(int64(node.Maxmem)),
			SslFingerprint: types.StringValue(node.SslFingerprint),
			Mem:            types.Int64Value(int64(node.Mem)),
			Id:             types.StringValue(node.Id),
			Node:           types.StringValue(node.Node),
			Disk:           types.Int64Value(node.Disk),
			Uptime:         types.Int64Value(int64(node.Uptime)),
			Maxdisk:        types.Int64Value(node.Maxdisk),
			Level:          types.StringValue(node.Level),
		}
		state.Nodes = append(state.Nodes, nodeState)
	}

	diags := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *nodeDataSource) Configure(_ context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}

	client, ok := request.ProviderData.(*proxmox.Client)
	if !ok {
		response.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *proxmox.Client, got %T. Please report this error to the developer", request.ProviderData),
		)
		return
	}

	d.client = client
}
