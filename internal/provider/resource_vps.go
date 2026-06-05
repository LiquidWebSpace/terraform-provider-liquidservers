package provider

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/liquidservers/terraform-provider-liquidservers/internal/liquidservers"
)

var (
	_ resource.Resource                = (*vpsResource)(nil)
	_ resource.ResourceWithConfigure   = (*vpsResource)(nil)
	_ resource.ResourceWithImportState = (*vpsResource)(nil)
)

type vpsResource struct {
	client *liquidservers.Client
}

func NewVPSResource() resource.Resource {
	return &vpsResource{}
}

func (r *vpsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vps"
}

func (r *vpsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a LiquidServers VPS through the Server Portal API.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Portal VM ID.",
			},
			"hostname": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Primary VPS hostname.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"label": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Friendly label stored in the portal.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"package_id": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "Optional portal package ID stored against the VPS.",
			},
			"client_reference": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Stable per-server automation reference. Set this explicitly when creating groups or importing/recovering state.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"os_template": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Portal OS template slug.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"node_name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Optional target node name. If omitted, the portal chooses its default node.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"location_name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Optional location label, inferred by the portal when omitted.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ipv4_count": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "Number of IPv4 addresses to allocate to this VPS.",
			},
			"ram_mb": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "Memory in MB.",
			},
			"disk_gb": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "Disk allocation in GB.",
			},
			"cpu_cores": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "Number of CPU cores.",
			},
			"bandwidth_gb": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "Bandwidth allocation in GB.",
			},
			"root_password": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "Optional initial root password. If omitted, the portal generates one.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Current portal status.",
			},
			"pending_action": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Current pending portal action, if any.",
			},
			"ip_address": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Primary IP address reported by the portal when available.",
			},
			"provider_vps_id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Underlying provider VPS ID.",
			},
			"last_synced_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Last portal sync timestamp.",
			},
		},
	}
}

func (r *vpsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*liquidservers.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected provider data", fmt.Sprintf("Expected *liquidservers.Client, got %T", req.ProviderData))
		return
	}

	r.client = client
}

func (r *vpsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan vpsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.ClientReference.IsNull() || plan.ClientReference.IsUnknown() {
		plan.ClientReference = types.StringValue("terraform-provider-liquidservers." + uuid.NewString())
	}

	input := liquidservers.CreateVPSRequest{
		Hostname:        plan.Hostname.ValueString(),
		Label:           plan.Label.ValueString(),
		PackageID:       packageIDPointer(plan.PackageID),
		ClientReference: plan.ClientReference.ValueString(),
		OSTemplate:      plan.OSTemplate.ValueString(),
		IPv4Count:       plan.IPv4Count.ValueInt64(),
		RAMMB:           plan.RAMMB.ValueInt64(),
		DiskGB:          plan.DiskGB.ValueInt64(),
		CPUCores:        plan.CPUCores.ValueInt64(),
		BandwidthGB:     plan.BandwidthGB.ValueInt64(),
	}
	if !plan.NodeName.IsNull() && !plan.NodeName.IsUnknown() {
		input.NodeName = plan.NodeName.ValueString()
	}
	if !plan.LocationName.IsNull() && !plan.LocationName.IsUnknown() {
		input.LocationName = plan.LocationName.ValueString()
	}
	if !plan.RootPassword.IsNull() && !plan.RootPassword.IsUnknown() {
		input.RootPassword = plan.RootPassword.ValueString()
		input.RootPasswordConfirmation = plan.RootPassword.ValueString()
	}

	idempotencyKey := "tf-create-" + plan.ClientReference.ValueString()
	vps, err := r.client.CreateVPS(ctx, input, idempotencyKey)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create VPS", err.Error())
		return
	}

	plan = populateVPSModel(plan, vps)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *vpsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state vpsModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := parseID(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid VPS ID", err.Error())
		return
	}

	vps, err := r.client.GetVPS(ctx, id, true)
	if err != nil {
		var apiErr liquidservers.APIError
		if errorAs(err, &apiErr) && apiErr.StatusCode == 422 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Unable to read VPS", err.Error())
		return
	}

	if vps.Status == "terminated" {
		resp.State.RemoveResource(ctx)
		return
	}

	state = populateVPSModel(state, vps)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *vpsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan vpsModel
	var state vpsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	for field, values := range map[string][2]int64{
		"ipv4_count":   {state.IPv4Count.ValueInt64(), plan.IPv4Count.ValueInt64()},
		"ram_mb":       {state.RAMMB.ValueInt64(), plan.RAMMB.ValueInt64()},
		"disk_gb":      {state.DiskGB.ValueInt64(), plan.DiskGB.ValueInt64()},
		"cpu_cores":    {state.CPUCores.ValueInt64(), plan.CPUCores.ValueInt64()},
		"bandwidth_gb": {state.BandwidthGB.ValueInt64(), plan.BandwidthGB.ValueInt64()},
	} {
		if values[1] < values[0] {
			resp.Diagnostics.AddAttributeError(path.Root(field), "Downgrade is not supported", "The LiquidServers API currently supports resource increases only. Increase the value or replace the VPS.")
			return
		}
	}

	id, err := parseID(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid VPS ID", err.Error())
		return
	}

	vps, err := r.client.UpgradeVPS(ctx, liquidservers.UpgradeVPSRequest{
		VMID:        id,
		PackageID:   packageIDPointer(plan.PackageID),
		IPv4Count:   plan.IPv4Count.ValueInt64(),
		RAMMB:       plan.RAMMB.ValueInt64(),
		DiskGB:      plan.DiskGB.ValueInt64(),
		CPUCores:    plan.CPUCores.ValueInt64(),
		BandwidthGB: plan.BandwidthGB.ValueInt64(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Unable to upgrade VPS", err.Error())
		return
	}

	plan = populateVPSModel(plan, vps)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *vpsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state vpsModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := parseID(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid VPS ID", err.Error())
		return
	}

	if err := r.client.TerminateVPS(ctx, id); err != nil {
		resp.Diagnostics.AddError("Unable to terminate VPS", err.Error())
		return
	}
}

func (r *vpsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func errorAs(err error, target any) bool {
	switch t := target.(type) {
	case *liquidservers.APIError:
		apiErr, ok := err.(liquidservers.APIError)
		if ok {
			*t = apiErr
			return true
		}
		apiErrPtr, ok := err.(*liquidservers.APIError)
		if ok {
			*t = *apiErrPtr
			return true
		}
	}

	return false
}

var _ = int64planmodifier.UseStateForUnknown
