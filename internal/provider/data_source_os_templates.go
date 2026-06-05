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
	_ datasource.DataSource              = (*osTemplatesDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*osTemplatesDataSource)(nil)
)

type osTemplatesDataSource struct {
	client *liquidservers.Client
}

type osTemplatesDataSourceModel struct {
	Templates []osTemplateModel `tfsdk:"templates"`
}

type osTemplateModel struct {
	ID   types.Int64  `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Slug types.String `tfsdk:"slug"`
	OSID types.Int64  `tfsdk:"osid"`
	Virt types.String `tfsdk:"virt"`
}

func NewOSTemplatesDataSource() datasource.DataSource {
	return &osTemplatesDataSource{}
}

func (d *osTemplatesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_os_templates"
}

func (d *osTemplatesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Lists available LiquidServers OS templates.",
		Attributes: map[string]schema.Attribute{
			"templates": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Available OS templates.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "Template ID.",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Template display name.",
						},
						"slug": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Template slug for VPS creation.",
						},
						"osid": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "Provider OS ID.",
						},
						"virt": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Virtualization type.",
						},
					},
				},
			},
		},
	}
}

func (d *osTemplatesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *osTemplatesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	templates, err := d.client.ListOSTemplates(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Unable to list OS templates", err.Error())
		return
	}

	state := osTemplatesDataSourceModel{
		Templates: make([]osTemplateModel, 0, len(templates)),
	}
	for _, template := range templates {
		state.Templates = append(state.Templates, osTemplateModel{
			ID:   types.Int64Value(template.ID),
			Name: types.StringValue(template.Name),
			Slug: types.StringValue(template.Slug),
			OSID: types.Int64Value(template.OSID),
			Virt: types.StringValue(template.Virt),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
