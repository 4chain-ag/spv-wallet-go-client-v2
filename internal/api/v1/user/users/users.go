package users

import (
	"context"
	"fmt"

	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/querybuilders"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/go-resty/resty/v2"
)

const route = "api/v1/users/current"

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

func (a *API) AccessKeys(ctx context.Context, transactionsOpts ...queries.AccessKeyQueryOption) (*queries.AccessKeyPage, error) {
	var query queries.AccessKeyQuery
	for _, o := range transactionsOpts {
		o(&query)
	}

	builderOpts := []querybuilders.QueryBuilderOption{
		querybuilders.WithMetadataFilter(query.Metadata),
		querybuilders.WithPageFilterQueryBuilder(query.PageFilter),
		querybuilders.WithFilterQueryBuilder(&accessKeyFilterBuilder{
			accessKeyFilter:    query.AccessKeyFilter,
			modelFilterBuilder: querybuilders.ModelFilterBuilder{ModelFilter: query.AccessKeyFilter.ModelFilter},
		}),
	}
	builder := querybuilders.NewQueryBuilder(builderOpts...)
	params, err := builder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build access keys query params: %w", err)
	}

	var result response.PageModel[response.AccessKey]
	URL := a.addr + "/keys"
	_, err = a.httpClient.
		R().
		SetContext(ctx).
		SetResult(&result).
		SetQueryParams(params.ParseToMap()).
		Get(URL)
	if err != nil {
		return nil, fmt.Errorf("HTTP response failure: %w", err)
	}

	return &result, nil
}

func (a *API) RevokeAccessKey(ctx context.Context, ID string) error {
	URL := a.addr + "/keys/" + ID
	_, err := a.httpClient.R().
		SetContext(ctx).
		Delete(URL)
	if err != nil {
		return fmt.Errorf("HTTP response failure: %w", err)
	}

	return nil
}

func NewAPI(addr string, httpClient *resty.Client) *API {
	return &API{
		addr:       addr + "/" + route,
		httpClient: httpClient,
	}
}
