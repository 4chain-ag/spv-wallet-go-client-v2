package configs

import (
	"context"
	"fmt"
	"net/url"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/errutil"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/go-resty/resty/v2"
)

const (
	userAPIRoute = "api/v1/configs"
	userAPI      = "User Shared Config API"
)

type UserAPI struct {
	url        *url.URL
	httpClient *resty.Client
}

func (u *UserAPI) SharedConfig(ctx context.Context) (*response.SharedConfig, error) {
	var result response.SharedConfig
	_, err := u.httpClient.
		R().
		SetContext(ctx).
		SetResult(&result).
		Get(u.url.JoinPath("shared").String())
	if err != nil {
		return nil, fmt.Errorf("HTTP response failure: %w", err)
	}

	return &result, nil
}

func NewUserAPI(url *url.URL, httpClient *resty.Client) *UserAPI {
	return &UserAPI{
		url:        url.JoinPath(userAPIRoute),
		httpClient: httpClient,
	}
}

func UserAPIErrorFormatter(action string, err error) *errutil.HTTPErrorFormatter {
	return &errutil.HTTPErrorFormatter{
		Action: action,
		API:    userAPI,
		Err:    err,
	}
}
