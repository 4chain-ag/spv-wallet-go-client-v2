package stats

import (
	"context"
	"fmt"
	"net/url"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/errutil"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/go-resty/resty/v2"
)

const (
	adminAPIRoute = "/v1/admin/stats"
	adminAPI      = "Admin Stats API"
)

type AdminAPI struct {
	url        *url.URL
	httpClient *resty.Client
}

func (a *AdminAPI) Stats(ctx context.Context) (*models.AdminStats, error) {
	var result models.AdminStats
	_, err := a.httpClient.
		R().
		SetContext(ctx).
		SetResult(&result).
		Get(a.url.String())
	if err != nil {
		return nil, fmt.Errorf("HTTP response failure: %w", err)
	}

	return &result, nil
}

func NewAdminAPI(url *url.URL, httpClient *resty.Client) *AdminAPI {
	return &AdminAPI{
		url:        url.JoinPath(adminAPIRoute),
		httpClient: httpClient,
	}
}

func AdminAPIErrorFormatter(action string, err error) *errutil.HTTPErrorFormatter {
	return &errutil.HTTPErrorFormatter{
		Action: action,
		API:    adminAPI,
		Err:    err,
	}
}
