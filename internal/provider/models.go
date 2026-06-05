package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

type vpsModel struct {
	ID              types.String `tfsdk:"id"`
	Hostname        types.String `tfsdk:"hostname"`
	Label           types.String `tfsdk:"label"`
	PackageID       types.Int64  `tfsdk:"package_id"`
	ClientReference types.String `tfsdk:"client_reference"`
	OSTemplate      types.String `tfsdk:"os_template"`
	NodeName        types.String `tfsdk:"node_name"`
	LocationName    types.String `tfsdk:"location_name"`
	IPv4Count       types.Int64  `tfsdk:"ipv4_count"`
	RAMMB           types.Int64  `tfsdk:"ram_mb"`
	DiskGB          types.Int64  `tfsdk:"disk_gb"`
	CPUCores        types.Int64  `tfsdk:"cpu_cores"`
	BandwidthGB     types.Int64  `tfsdk:"bandwidth_gb"`
	RootPassword    types.String `tfsdk:"root_password"`
	Status          types.String `tfsdk:"status"`
	PendingAction   types.String `tfsdk:"pending_action"`
	IPAddress       types.String `tfsdk:"ip_address"`
	ProviderVPSID   types.Int64  `tfsdk:"provider_vps_id"`
	LastSyncedAt    types.String `tfsdk:"last_synced_at"`
}
