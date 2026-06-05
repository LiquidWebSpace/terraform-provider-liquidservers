package liquidservers

import (
	"context"
	"fmt"
	"net/url"
)

type VPS struct {
	ID            int64   `json:"id"`
	OwnerUserID   int64   `json:"owner_user_id"`
	ProviderVPSID int64   `json:"virtualizor_vps_id"`
	Hostname      string  `json:"hostname"`
	Label         string  `json:"label"`
	PackageID     *int64  `json:"package_id"`
	ClientRef     string  `json:"client_reference"`
	NodeName      string  `json:"node_name"`
	LocationName  string  `json:"location_name"`
	OSTemplate    string  `json:"os_template"`
	IPv4Count     int64   `json:"ipv4_count"`
	RAMMB         int64   `json:"ram_mb"`
	DiskGB        int64   `json:"disk_gb"`
	CPUCores      int64   `json:"cpu_cores"`
	BandwidthGB   int64   `json:"bandwidth_gb"`
	Status        string  `json:"status"`
	PendingAction *string `json:"pending_action"`
	IPAddress     string  `json:"ip_address"`
	LastSyncedAt  string  `json:"last_synced_at"`
	TerminatedAt  *string `json:"terminated_at"`
}

type CreateVPSRequest struct {
	Hostname                 string `json:"hostname"`
	Label                    string `json:"label"`
	PackageID                *int64 `json:"package_id,omitempty"`
	ClientReference          string `json:"client_reference,omitempty"`
	OSTemplate               string `json:"os_template"`
	NodeName                 string `json:"node_name,omitempty"`
	LocationName             string `json:"location_name,omitempty"`
	IPv4Count                int64  `json:"ipv4_count"`
	RAMMB                    int64  `json:"ram_mb"`
	DiskGB                   int64  `json:"disk_gb"`
	CPUCores                 int64  `json:"cpu_cores"`
	BandwidthGB              int64  `json:"bandwidth_gb"`
	RootPassword             string `json:"root_password,omitempty"`
	RootPasswordConfirmation string `json:"root_password_confirmation,omitempty"`
}

type UpgradeVPSRequest struct {
	VMID        int64  `json:"vm_id"`
	PackageID   *int64 `json:"package_id,omitempty"`
	IPv4Count   int64  `json:"ipv4_count"`
	RAMMB       int64  `json:"ram_mb"`
	DiskGB      int64  `json:"disk_gb"`
	CPUCores    int64  `json:"cpu_cores"`
	BandwidthGB int64  `json:"bandwidth_gb"`
}

func (c *Client) CreateVPS(ctx context.Context, input CreateVPSRequest, idempotencyKey string) (*VPS, error) {
	var response struct {
		VMID int64 `json:"vm_id"`
		Data *VPS  `json:"data"`
	}
	if err := c.do(ctx, httpPost, "/api/vps", input, idempotencyKey, &response); err != nil {
		return nil, err
	}
	if response.Data != nil && response.Data.ID > 0 {
		return response.Data, nil
	}
	if response.VMID <= 0 {
		return nil, fmt.Errorf("create response did not include vm_id")
	}

	return c.GetVPS(ctx, response.VMID, true)
}

func (c *Client) GetVPS(ctx context.Context, id int64, includeTerminated bool) (*VPS, error) {
	path := fmt.Sprintf("/api/vps?id=%d", id)
	if includeTerminated {
		path += "&include_terminated=1"
	}

	var response struct {
		Data VPS `json:"data"`
	}
	if err := c.do(ctx, httpGet, path, nil, "", &response); err != nil {
		return nil, err
	}
	if response.Data.ID == 0 {
		return nil, fmt.Errorf("VPS %d not found", id)
	}

	return &response.Data, nil
}

func (c *Client) ListVPS(ctx context.Context) ([]VPS, error) {
	var response struct {
		Data []VPS `json:"data"`
	}
	if err := c.do(ctx, httpGet, "/api/vps", nil, "", &response); err != nil {
		return nil, err
	}

	return response.Data, nil
}

func (c *Client) UpgradeVPS(ctx context.Context, input UpgradeVPSRequest) (*VPS, error) {
	if err := c.do(ctx, httpPost, "/api/vps/upgrade", input, "", nil); err != nil {
		return nil, err
	}

	return c.GetVPS(ctx, input.VMID, true)
}

func (c *Client) TerminateVPS(ctx context.Context, id int64) error {
	body := map[string]any{
		"vm_id":  id,
		"action": "terminate",
	}

	return c.do(ctx, httpPost, "/api/vps/action", body, "", nil)
}

func (c *Client) FindVPSByClientReference(ctx context.Context, clientReference string) (*VPS, error) {
	if clientReference == "" {
		return nil, nil
	}

	vpsList, err := c.ListVPS(ctx)
	if err != nil {
		return nil, err
	}
	for _, item := range vpsList {
		if item.ClientRef == clientReference {
			return &item, nil
		}
	}

	return nil, nil
}

type OSTemplate struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
	OSID int64  `json:"osid"`
	Virt string `json:"virt"`
}

func (c *Client) ListOSTemplates(ctx context.Context) ([]OSTemplate, error) {
	var response struct {
		Data []OSTemplate `json:"data"`
	}
	if err := c.do(ctx, httpGet, "/api/os", nil, "", &response); err != nil {
		return nil, err
	}

	return response.Data, nil
}

const (
	httpGet  = "GET"
	httpPost = "POST"
)

func QueryEscape(value string) string {
	return url.QueryEscape(value)
}
