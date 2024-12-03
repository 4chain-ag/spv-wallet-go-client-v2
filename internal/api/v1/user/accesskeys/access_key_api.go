package accesskeys

import (
	"context"
	"fmt"
	"net/url"

	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/errutil"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/querybuilders"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/go-resty/resty/v2"
)

const route = "api/v1/users/current"

type API struct {
	url        *url.URL
	httpClient *resty.Client
}

func (a *API) GenerateAccessKey(ctx context.Context, cmd *commands.GenerateAccessKey) (*response.AccessKey, error) {
	var result response.AccessKey

	_, err := a.httpClient.R().
		SetContext(ctx).
		SetResult(&result).
		SetBody(cmd).
		Post(a.url.JoinPath("keys").String())
	if err != nil {
		return nil, fmt.Errorf("HTTP response failure: %w", err)
	}

	return &result, nil
}

func (a *API) AccessKey(ctx context.Context, ID string) (*response.AccessKey, error) {
	var result response.AccessKey

	_, err := a.httpClient.R().
		SetContext(ctx).
		SetResult(&result).
		Get(a.url.JoinPath("keys", ID).String())
	if err != nil {
		return nil, fmt.Errorf("HTTP response failure: %w", err)
	}

	return &result, nil
}

func (a *API) AccessKeys(ctx context.Context, opts ...queries.AccessKeyQueryOption) (*queries.AccessKeyPage, error) {
	var query queries.AccessKeyQuery
	for _, o := range opts {
		o(&query)
	}

	queryBuilder := querybuilders.NewQueryBuilder(
		querybuilders.WithMetadataFilter(query.Metadata),
		querybuilders.WithPageFilter(query.PageFilter),
		querybuilders.WithFilterQueryBuilder(&AccessKeyFilterQueryBuilder{
			AccessKeyFilter:    query.AccessKeyFilter,
			ModelFilterBuilder: querybuilders.ModelFilterBuilder{ModelFilter: query.AccessKeyFilter.ModelFilter},
		}),
	)
	params, err := queryBuilder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build access keys query params: %w", err)
	}

	var result response.PageModel[response.AccessKey]
	_, err = a.httpClient.
		R().
		SetContext(ctx).
		SetResult(&result).
		SetQueryParams(params.ParseToMap()).
		Get(a.url.JoinPath("keys").String())
	if err != nil {
		return nil, fmt.Errorf("HTTP response failure: %w", err)
	}

	return &result, nil
}

func (a *API) RevokeAccessKey(ctx context.Context, ID string) error {
	_, err := a.httpClient.R().
		SetContext(ctx).
		Delete(a.url.JoinPath("keys", ID).String())
	if err != nil {
		return fmt.Errorf("HTTP response failure: %w", err)
	}

	return nil
}

func NewAPI(url *url.URL, httpClient *resty.Client) *API {
	return &API{
		url:        url.JoinPath(route),
		httpClient: httpClient,
	}
}

func HTTPErrorFormatter(action string, err error) *errutil.HTTPErrorFormatter {
	return &errutil.HTTPErrorFormatter{
		Action: action,
		API:    "User Access Keys API",
		Err:    err,
	}
}