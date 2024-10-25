package configurations

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bitcoin-sv/spv-wallet/models/response"
)

// HTTP is an interface that defines methods for making HTTP requests.
type HTTP interface {
	Get(ctx context.Context, path string) ([]byte, error)
}

// API represents an API client with a specific base path and an HTTP interface.
type API struct {
	Path string
	HTTP HTTP
}

// SharedConfig fetches the shared configuration from the SPV Wallet API.
// It constructs the request path and unmarshals the response into a SharedConfig struct.
func (a *API) SharedConfig(ctx context.Context) (*response.SharedConfig, error) {
	path := a.Path + "/shared"
	body, err := a.HTTP.Get(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to get shared config from %s: %w", path, err)
	}

	var dst response.SharedConfig
	if err := json.Unmarshal(body, &dst); err != nil {
		return nil, fmt.Errorf("failed to unmarshal HTTP response into SharedConfig: %w", err)
	}
	return &dst, nil
}
