package spvwallet

import (
	"context"
	"fmt"
	"net/url"

	bip32 "github.com/bitcoin-sv/go-sdk/compat/bip32"
	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/config"
	xpubs "github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/admin/users"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/auth"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/restyutil"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

// AdminAPI provides a convenient interface for interacting with admin-related APIs.
// It abstracts the complexities of making HTTP requests and handling responses, allowing
// developers to interact with the admin API endpoints in a simplified and consistent manner.
type AdminAPI struct {
	xpubsAPI *xpubs.API // Internal API for managing operations related to XPubs.
}

// CreateXPub creates a new XPub record using the Admin XPubs API.
// The provided command contains the necessary parameters to define the XPub record.
//
// The API response is unmarshaled into a *response.Xpub struct. If the request
// encounters an error or the response cannot be properly decoded, an error is returned.
func (a *AdminAPI) CreateXPub(ctx context.Context, cmd *commands.CreateUserXpub) (*response.Xpub, error) {
	res, err := a.xpubsAPI.CreateXPub(ctx, cmd)
	if err != nil {
		return nil, xpubs.HTTPErrorFormatter("failed to create XPub", err).FormatPostErr()
	}
	return res, nil
}

// XPubs retrieves a paginated list of user XPubs from the Admin XPubs API.
// The response includes user XPubs along with pagination metadata, such as
// the current page number, sort order, and the field used for sorting (sortBy).
//
// Query parameters can be configured using optional query options. These options allow
// filtering based on metadata, pagination settings, or specific XPub attributes.
//
// The API response is unmarshaled into a *queries.XPubPage struct. If the API call fails
// or the response cannot be decoded, an error is returned.
func (a *AdminAPI) XPubs(ctx context.Context, opts ...queries.XPubQueryOption) (*queries.XPubPage, error) {
	res, err := a.xpubsAPI.XPubs(ctx, opts...)
	if err != nil {
		return nil, xpubs.HTTPErrorFormatter("failed to retrieve XPubs page", err).FormatGetErr()
	}

	return res, nil
}

// NewAdminAPI initializes a new instance of AdminAPI.
// It configures the API client using the provided configuration and xPub key for authentication.
// If any step fails, an appropriate error is returned.
func NewAdminAPI(cfg config.Config, xPub string) (*AdminAPI, error) {
	url, err := url.Parse(cfg.Addr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse address as url.URL: %w", err)
	}

	key, err := bip32.GetHDKeyFromExtendedPublicKey(xPub)
	if err != nil {
		return nil, fmt.Errorf("failed to generate HD key from xPub: %w", err)
	}

	authenticator, err := auth.NewXpubOnlyAuthenticator(key)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize xPub authenticator: %w", err)
	}

	httpClient := restyutil.NewHTTPClient(cfg, authenticator)
	return &AdminAPI{xpubsAPI: xpubs.NewAPI(url, httpClient)}, nil
}
