package accesskeys

import (
	"context"
	"fmt"
	"net/url"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/accesskeys/builders"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/errutil"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/querybuilders"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/go-resty/resty/v2"
)

const (
	adminAPIRoute = "/api/v1/admin/users/keys"
	adminAPI      = "Admin Access Keys API"
)

type AdminAPI struct {
	httpClient *resty.Client
	url        *url.URL
}

func (a *AdminAPI) AccessKeys(ctx context.Context, opts ...queries.AdminAccessKeyQueryOption) (*queries.AccessKeyPage, error) {
	var query queries.AdminAccessKeyQuery
	for _, o := range opts {
		o(&query)
	}

	queryBuilder := querybuilders.NewQueryBuilder(
		querybuilders.WithMetadataFilter(query.Metadata),
		querybuilders.WithPageFilter(query.PageFilter),
		querybuilders.WithFilterQueryBuilder(&builders.AdminAccessKeyFilterQueryBuilder{AccessKeyFilter: query.AdminAccessKeyFilter}),
	)
	params, err := queryBuilder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build access keys query params: %w", err)
	}

	var result queries.AccessKeyPage
	_, err = a.httpClient.
		R().
		SetContext(ctx).
		SetResult(&result).
		SetQueryParams(params.ParseToMap()).
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
		API:    adminAPIRoute,
		Err:    err,
	}
}
