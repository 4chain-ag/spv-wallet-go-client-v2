package invitations

import (
	"context"
	"fmt"
	"net/url"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/errutil"
	"github.com/go-resty/resty/v2"
)

const (
	adminAPIRoute = "api/v1/admin/invitations"
	adminAPI      = "Admin Invitations API"
)

type AdminAPI struct {
	httpClient *resty.Client
	url        *url.URL
}

func (a *AdminAPI) AcceptInvitation(ctx context.Context, ID string) error {
	URL := a.url.JoinPath(ID).String()
	_, err := a.httpClient.
		R().
		SetContext(ctx).
		Post(URL)
	if err != nil {
		return fmt.Errorf("HTTP response failure: %w", err)
	}

	return nil
}

func (a *AdminAPI) RejectInvitation(ctx context.Context, ID string) error {
	URL := a.url.JoinPath(ID).String()
	_, err := a.httpClient.
		R().
		SetContext(ctx).
		Delete(URL)
	if err != nil {
		return fmt.Errorf("HTTP response failure: %w", err)
	}

	return nil
}

func NewAdminAPI(url *url.URL, httpClient *resty.Client) *AdminAPI {
	return &AdminAPI{url: url.JoinPath(adminAPIRoute), httpClient: httpClient}
}

func AdminAPIErrorFormatter(action string, err error) *errutil.HTTPErrorFormatter {
	return &errutil.HTTPErrorFormatter{
		Action: action,
		API:    adminAPI,
		Err:    err,
	}
}
