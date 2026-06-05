package liquidservers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	endpoint   string
	apiKey     string
	httpClient *http.Client
}

type APIError struct {
	StatusCode int
	Message    string
}

func (e APIError) Error() string {
	if e.Message == "" {
		return fmt.Sprintf("LiquidServers API request failed with HTTP %d", e.StatusCode)
	}

	return fmt.Sprintf("LiquidServers API request failed with HTTP %d: %s", e.StatusCode, e.Message)
}

func NewClient(endpoint, apiKey string) (*Client, error) {
	endpoint = strings.TrimRight(strings.TrimSpace(endpoint), "/")
	if endpoint == "" {
		return nil, fmt.Errorf("endpoint is required")
	}
	if _, err := url.ParseRequestURI(endpoint); err != nil {
		return nil, fmt.Errorf("endpoint must be a valid URL: %w", err)
	}
	if strings.TrimSpace(apiKey) == "" {
		return nil, fmt.Errorf("api_key is required")
	}

	return &Client{
		endpoint: endpoint,
		apiKey:   strings.TrimSpace(apiKey),
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
	}, nil
}

func (c *Client) do(ctx context.Context, method, path string, body any, idempotencyKey string, out any) error {
	var reader io.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reader = bytes.NewReader(payload)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.endpoint+path, reader)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("X-API-Key", c.apiKey)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if idempotencyKey != "" {
		req.Header.Set("Idempotency-Key", idempotencyKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		var errorPayload struct {
			Error string `json:"error"`
		}
		_ = json.Unmarshal(raw, &errorPayload)
		return APIError{StatusCode: resp.StatusCode, Message: firstNonEmpty(errorPayload.Error, string(raw))}
	}

	if out == nil || len(raw) == 0 {
		return nil
	}

	return json.Unmarshal(raw, out)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			return value
		}
	}

	return ""
}
