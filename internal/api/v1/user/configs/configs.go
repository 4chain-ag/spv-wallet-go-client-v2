// Package configs provides functionality to communicate with the SPV Wallet API
// endpoints related to user configuration.
package configs

import (
	"context"
	"fmt"

	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/go-resty/resty/v2"
)

const route = "api/v1/configs"

// API represents a client for accessing user configuration API endpoints.
type API struct {
	addr string        // Base address for the configuration endpoints.
	cli  *resty.Client // HTTP client used for making request.
}

// SharedConfig fetches the shared configuration from the SPV Wallet API.
// It constructs the request path and unmarshals the response into a SharedConfig struct.
func (a *API) SharedConfig(ctx context.Context) (*response.SharedConfig, error) {
	var spvErr models.SPVError
	var result response.SharedConfig

	resp, err := a.cli.
		R().
		SetContext(ctx).
		SetResult(&result).
		SetError(&spvErr).
		Get(a.addr + "/shared")
	if err != nil {
		return nil, fmt.Errorf("HTTP error: %w", err)
	}
	if resp.IsError() {
		return nil, spvErr
	}
	return &result, nil
}

// NewAPI constructs a new instance of the configurations API client (`API`).
// This client provides functionality for making HTTP requests to endpoints
// within the user configurations domain. The client uses the provided SPV Wallet
// API address as the base URL and relies on a `resty.Client` for handling HTTP requests.
func NewAPI(addr string, cli *resty.Client) *API {
	return &API{
		addr: addr + "/" + route,
		cli:  cli,
	}
}
