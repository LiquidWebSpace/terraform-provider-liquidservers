package provider

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/liquidservers/terraform-provider-liquidservers/internal/liquidservers"
)

func packageIDPointer(value types.Int64) *int64 {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}

	v := value.ValueInt64()
	return &v
}

func stringOrNull(value string) types.String {
	if value == "" {
		return types.StringNull()
	}

	return types.StringValue(value)
}

func pendingActionOrNull(value *string) types.String {
	if value == nil || *value == "" {
		return types.StringNull()
	}

	return types.StringValue(*value)
}

func populateVPSModel(plan vpsModel, vps *liquidservers.VPS) vpsModel {
	plan.ID = types.StringValue(strconv.FormatInt(vps.ID, 10))
	plan.Hostname = types.StringValue(vps.Hostname)
	plan.Label = types.StringValue(vps.Label)
	if vps.PackageID == nil {
		plan.PackageID = types.Int64Null()
	} else {
		plan.PackageID = types.Int64Value(*vps.PackageID)
	}
	plan.ClientReference = stringOrNull(vps.ClientRef)
	plan.OSTemplate = types.StringValue(vps.OSTemplate)
	plan.NodeName = stringOrNull(vps.NodeName)
	plan.LocationName = stringOrNull(vps.LocationName)
	plan.IPv4Count = types.Int64Value(vps.IPv4Count)
	plan.RAMMB = types.Int64Value(vps.RAMMB)
	plan.DiskGB = types.Int64Value(vps.DiskGB)
	plan.CPUCores = types.Int64Value(vps.CPUCores)
	plan.BandwidthGB = types.Int64Value(vps.BandwidthGB)
	plan.Status = types.StringValue(vps.Status)
	plan.PendingAction = pendingActionOrNull(vps.PendingAction)
	plan.IPAddress = stringOrNull(vps.IPAddress)
	plan.ProviderVPSID = types.Int64Value(vps.ProviderVPSID)
	plan.LastSyncedAt = stringOrNull(vps.LastSyncedAt)

	return plan
}

func parseID(value string) (int64, error) {
	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("expected positive numeric ID, got %q", value)
	}

	return id, nil
}
