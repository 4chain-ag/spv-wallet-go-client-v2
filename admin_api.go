package client

import (
	"context"
	"fmt"
	"net/url"

	bip32 "github.com/bitcoin-sv/go-sdk/compat/bip32"
	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	xpubs "github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/admin/users"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/auth"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

type AdminAPI struct {
	xpubsAPI *xpubs.API
}

func (a *AdminAPI) CreateXPub(ctx context.Context, cmd *commands.CreateUserXpub) (*response.Xpub, error) {
	res, err := a.xpubsAPI.CreateXPub(ctx, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to create xpub record by calling the admin xpubs API: %w", err)
	}

	return res, nil
}

func (a *AdminAPI) XPubs(ctx context.Context, opts ...queries.XPubQueryOption) (*queries.XPubPage, error) {
	res, err := a.xpubsAPI.XPubs(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create xpub record by calling the admin xpubs API: %w", err)
	}

	return res, nil
}

func NewAdminAPI(cfg Config, xPub string) (*AdminAPI, error) {
	url, err := url.Parse(cfg.Addr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse addr to url.URL: %w", err)
	}

	key, err := bip32.GetHDKeyFromExtendedPublicKey(xPub)
	if err != nil {
		return nil, fmt.Errorf("failed to generate HD key from xPub: %w", err)
	}

	authenticator, err := auth.NewXpubOnlyAuthenticator(key)
	if err != nil {
		return nil, fmt.Errorf("failed to intialized xpub authenticator: %w", err)
	}

	httpClient := newRestyClient(cfg, authenticator)
	return &AdminAPI{xpubsAPI: xpubs.NewAPI(url, httpClient)}, nil
}
