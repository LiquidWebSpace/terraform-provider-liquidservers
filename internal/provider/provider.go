package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/liquidservers/terraform-provider-liquidservers/internal/liquidservers"
)

var _ provider.Provider = (*liquidServersProvider)(nil)

type liquidServersProvider struct {
	version string
}

type liquidServersProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	APIKey   types.String `tfsdk:"api_key"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &liquidServersProvider{version: version}
	}
}

func (p *liquidServersProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "liquidservers"
	resp.Version = p.version
}

func (p *liquidServersProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manage LiquidServers VPS resources through the LiquidServers Server Portal API.",
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "LiquidServers Server Portal base URL. Can also be set with `LIQUIDSERVERS_ENDPOINT`.",
			},
			"api_key": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "LiquidServers API key. Can also be set with `LIQUIDSERVERS_API_KEY`.",
			},
		},
	}
}

func (p *liquidServersProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config liquidServersProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := os.Getenv("LIQUIDSERVERS_ENDPOINT")
	if !config.Endpoint.IsNull() && !config.Endpoint.IsUnknown() {
		endpoint = config.Endpoint.ValueString()
	}

	apiKey := os.Getenv("LIQUIDSERVERS_API_KEY")
	if !config.APIKey.IsNull() && !config.APIKey.IsUnknown() {
		apiKey = config.APIKey.ValueString()
	}

	if endpoint == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Missing LiquidServers endpoint",
			"Set endpoint in the provider block or LIQUIDSERVERS_ENDPOINT in the environment.",
		)
	}

	if apiKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing LiquidServers API key",
			"Set api_key in the provider block or LIQUIDSERVERS_API_KEY in the environment.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	client, err := liquidservers.NewClient(endpoint, apiKey)
	if err != nil {
		resp.Diagnostics.AddError("Invalid LiquidServers provider configuration", err.Error())
		return
	}

	resp.ResourceData = client
	resp.DataSourceData = client
}

func (p *liquidServersProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewVPSResource,
	}
}

func (p *liquidServersProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewOSTemplatesDataSource,
		NewVPSDataSource,
	}
}
