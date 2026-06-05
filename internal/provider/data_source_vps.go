package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/liquidservers/terraform-provider-liquidservers/internal/liquidservers"
)

var (
	_ datasource.DataSource              = (*vpsDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*vpsDataSource)(nil)
)

type vpsDataSource struct {
	client *liquidservers.Client
}

func NewVPSDataSource() datasource.DataSource {
	return &vpsDataSource{}
}

func (d *vpsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vps"
}

func (d *vpsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Reads an existing LiquidServers VPS by portal VM ID.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Portal VM ID.",
			},
			"hostname": schema.StringAttribute{
				Computed: true,
			},
			"label": schema.StringAttribute{
				Computed: true,
			},
			"package_id": schema.Int64Attribute{
				Computed: true,
			},
			"client_reference": schema.StringAttribute{
				Computed: true,
			},
			"os_template": schema.StringAttribute{
				Computed: true,
			},
			"node_name": schema.StringAttribute{
				Computed: true,
			},
			"location_name": schema.StringAttribute{
				Computed: true,
			},
			"ipv4_count": schema.Int64Attribute{
				Computed: true,
			},
			"ram_mb": schema.Int64Attribute{
				Computed: true,
			},
			"disk_gb": schema.Int64Attribute{
				Computed: true,
			},
			"cpu_cores": schema.Int64Attribute{
				Computed: true,
			},
			"bandwidth_gb": schema.Int64Attribute{
				Computed: true,
			},
			"root_password": schema.StringAttribute{
				Computed:  true,
				Sensitive: true,
			},
			"status": schema.StringAttribute{
				Computed: true,
			},
			"pending_action": schema.StringAttribute{
				Computed: true,
			},
			"ip_address": schema.StringAttribute{
				Computed: true,
			},
			"provider_vps_id": schema.Int64Attribute{
				Computed: true,
			},
			"last_synced_at": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d *vpsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*liquidservers.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected provider data", fmt.Sprintf("Expected *liquidservers.Client, got %T", req.ProviderData))
		return
	}

	d.client = client
}

func (d *vpsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state vpsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := parseID(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid VPS ID", err.Error())
		return
	}

	vps, err := d.client.GetVPS(ctx, id, true)
	if err != nil {
		resp.Diagnostics.AddError("Unable to read VPS", err.Error())
		return
	}

	state.RootPassword = types.StringNull()
	state = populateVPSModel(state, vps)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
