package users

import (
	"context"
	"fmt"

	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/go-resty/resty/v2"
)

const route = "api/v1/users/current"

// TODO: 1. Implement search access key endpoint call.

type API struct {
	addr       string
	httpClient *resty.Client
}

func (a *API) XPub(ctx context.Context) (*response.Xpub, error) {
	var result response.Xpub
	_, err := a.httpClient.
		R().
		SetContext(ctx).
		SetResult(&result).
		Get(a.addr)
	if err != nil {
		return nil, fmt.Errorf("HTTP response failure: %w", err)
	}

	return &result, nil
}

func (a *API) UpdateXPubMetadata(ctx context.Context, cmd *commands.UpdateXPubMetadata) (*response.Xpub, error) {
	var result response.Xpub
	_, err := a.httpClient.R().
		SetContext(ctx).
		SetResult(&result).
		SetBody(cmd).
		Patch(a.addr)
	if err != nil {
		return nil, fmt.Errorf("HTTP response failure: %w", err)
	}

	return &result, nil
}

func (a *API) GenerateAccessKey(ctx context.Context, cmd *commands.GenerateAccessKey) (*response.AccessKey, error) {
	var result response.AccessKey

	URL := a.addr + "/keys"
	_, err := a.httpClient.R().
		SetContext(ctx).
		SetResult(&result).
		SetBody(cmd).
		Post(URL)
	if err != nil {
		return nil, fmt.Errorf("HTTP response failure: %w", err)
	}

	return &result, nil
}

func (a *API) AccessKey(ctx context.Context, ID string) (*response.AccessKey, error) {
	var result response.AccessKey

	URL := a.addr + "/keys/" + ID
	_, err := a.httpClient.R().
		SetContext(ctx).
		SetResult(&result).
		Get(URL)
	if err != nil {
		return nil, fmt.Errorf("HTTP response failure: %w", err)
	}

	return &result, nil
}

func (a *API) RevokeAccessKey(ctx context.Context, ID string) (*response.AccessKey, error) {
	var result response.AccessKey

	URL := a.addr + "/keys/" + ID
	_, err := a.httpClient.R().
		SetContext(ctx).
		SetResult(&result).
		Delete(URL)
	if err != nil {
		return nil, fmt.Errorf("HTTP response failure: %w", err)
	}

	return &result, nil
}

func NewAPI(addr string, httpClient *resty.Client) *API {
	return &API{
		addr:       addr + "/" + route,
		httpClient: httpClient,
	}
}
