package status

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/errutil"
	"github.com/go-resty/resty/v2"
)

const (
	adminAPIRoute = "v1/admin/status"
	adminAPI      = "Admin Status API"
)

type AdminAPI struct {
	httpClient *resty.Client
	url        *url.URL
}

func (a *AdminAPI) Status(ctx context.Context) (bool, error) {
	res, err := a.httpClient.
		R().
		SetContext(ctx).
		Get(a.url.String())
	if err != nil {
		if res.StatusCode() == http.StatusUnauthorized {
			return false, nil
		}
		return false, fmt.Errorf("HTTP response failure: %w", err)
	}

	return true, nil
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
