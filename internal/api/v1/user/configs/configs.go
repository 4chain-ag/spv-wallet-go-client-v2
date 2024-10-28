package configs

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bitcoin-sv/spv-wallet/models/response"
)

// route defines the API endpoint path for accessing configurations in the SPV Wallet API.
// This constant is appended to the base URL to form the complete path for configuration
// requests.
const route = "api/v1/configs"

// HTTPGetter is an interface that defines methods for making HTTP GET request.
type HTTPGetter interface {
	Get(ctx context.Context, path string) ([]byte, error)
}

// API represents a client for interacting with an API at a specific base path,
// using an HTTPGetter interface for making requests.
type API struct {
	addr string     // addr defines the base URL for the user configurations API.
	cli  HTTPGetter // cli is used for intializing the external SPV Wallet API HTTP calls.
}

// SharedConfig fetches the shared configuration from the SPV Wallet API.
// It constructs the request path and unmarshals the response into a SharedConfig struct.
func (a *API) SharedConfig(ctx context.Context) (*response.SharedConfig, error) {
	path := a.addr + "/shared"
	body, err := a.cli.Get(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to get shared config from %s: %w", path, err)
	}

	var dst response.SharedConfig
	if err := json.Unmarshal(body, &dst); err != nil {
		return nil, fmt.Errorf("failed to unmarshal HTTP response into SharedConfig: %w", err)
	}
	return &dst, nil
}

// NewAPI constructs a new instance of the configurations API client (`API`) using
// the provided SPV Wallet API address and HTTPGetter implementation. This client is
// responsible for constructing HTTP requests and making calls to the appropriate
// endpoints within the user configurations domain.
func NewAPI(addr string, h HTTPGetter) *API {
	return &API{
		addr: addr + "/" + route,
		cli:  h,
	}
}
