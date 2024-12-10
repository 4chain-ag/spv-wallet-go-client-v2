package invitations

import (
	"context"
	"fmt"
	"net/url"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/errutil"
	"github.com/go-resty/resty/v2"
)

const (
	userAPIRoute = "api/v1/invitations"
	userAPI      = "User Invitations API"
)

type UserAPI struct {
	url        *url.URL
	httpClient *resty.Client
}

func (u *UserAPI) AcceptInvitation(ctx context.Context, paymail string) error {
	_, err := u.httpClient.
		R().
		SetContext(ctx).
		Post(u.url.JoinPath(paymail, "contacts").String())
	if err != nil {
		return fmt.Errorf("HTTP response failure: %w", err)
	}

	return nil
}

func (u *UserAPI) RejectInvitation(ctx context.Context, paymail string) error {
	_, err := u.httpClient.
		R().
		SetContext(ctx).
		Delete(u.url.JoinPath(paymail).String())
	if err != nil {
		return fmt.Errorf("HTTP response failure: %w", err)
	}

	return nil
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
