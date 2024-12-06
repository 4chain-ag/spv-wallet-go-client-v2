package utxos

import (
	"context"
	"fmt"
	"net/url"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/errutil"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/querybuilders"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/go-resty/resty/v2"
)

const (
	route = "/api/v1/admin/utxos"
	api   = "Admin UTXOs API"
)

type API struct {
	httpClient *resty.Client
	url        *url.URL
}

func (a *API) UTXOs(ctx context.Context, opts ...queries.AdminUtxoQueryOption) (*queries.UtxosPage, error) {
	var query queries.AdminUtxoQuery
	for _, o := range opts {
		o(&query)
	}

	queryBuilder := querybuilders.NewQueryBuilder(
		querybuilders.WithMetadataFilter(query.Metadata),
		querybuilders.WithPageFilter(query.PageFilter),
		querybuilders.WithFilterQueryBuilder(&adminUtxoFilterQueryBuilder{utxoFilter: query.UtxoFilter}),
	)
	params, err := queryBuilder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build utxos query params: %w", err)
	}

	var result queries.UtxosPage
	_, err = a.httpClient.
		R().
		SetContext(ctx).
		SetQueryParams(params.ParseToMap()).
		SetResult(&result).
		Get(a.url.String())
	if err != nil {
		return nil, fmt.Errorf("HTTP response failure: %w", err)
	}

	return &result, nil
}

func NewAPI(url *url.URL, httpClient *resty.Client) *API {
	return &API{url: url.JoinPath(route), httpClient: httpClient}
}

func HTTPErrorFormatter(action string, err error) *errutil.HTTPErrorFormatter {
	return &errutil.HTTPErrorFormatter{
		Action: action,
		API:    api,
		Err:    err,
	}
}