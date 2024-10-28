package httpx

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/auth"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/go-resty/resty/v2"
)

// Resty provides a simplified facade over the resty.Client for making HTTP requests with
// contextual authentication headers. This facade is designed to handle client configuration,
// setup, and authentication, abstracting away the underlying complexities of the resty library.
type Resty struct {
	cli *resty.Client
}

// Get performs an HTTP GET request to the specified path using the configured resty.Client.
// It attaches authentication headers and handles errors gracefully, returning the response body
// or an error if the request fails.
func (r *Resty) Get(ctx context.Context, path string) ([]byte, error) {
	req := r.cli.R().SetContext(ctx)
	res, err := req.Get(path)
	if err != nil {
		return nil, fmt.Errorf("HTTP GET %s: %w", r.cli.BaseURL+path, err)
	}
	return res.Body(), nil
}

// NewResty creates and configures a Resty instance for authenticated requests.
func NewResty(addr string, cfg *auth.HeaderConfig) (*Resty, error) {
	b, err := auth.NewHeaderBuilder(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to build HTTP auth header: %w", err)

	}
	cli := resty.New().
		SetError(&models.SPVError{}).
		SetBaseURL(addr).
		OnBeforeRequest(authHeaderMiddleware(b))

	return &Resty{cli: cli}, nil
}

// AuthHeaderMiddleware creates and attaches authentication headers to requests.
func authHeaderMiddleware(b *auth.HeaderBuilder) func(c *resty.Client, r *resty.Request) error {
	return func(_ *resty.Client, r *resty.Request) error {
		switch r.Method {
		case http.MethodGet:
			auth, err := b.BuildWithoutBody() // Constructs an HTTP header.
			if err != nil {
				return fmt.Errorf("failed to build HTTP auth headers: %w", err)
			}
			r.SetHeaderMultiValues(auth)
		}
		return nil
	}
}
