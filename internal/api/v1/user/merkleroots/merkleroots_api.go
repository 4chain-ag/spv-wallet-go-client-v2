package merkleroots

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	customerr "github.com/bitcoin-sv/spv-wallet-go-client/errors"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/querybuilders"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/go-resty/resty/v2"
)

const route = "api/v1/merkleroots"

type DB interface {
	GetLastMerkleRoot() string
	SaveMerkleRoots([]models.MerkleRoot) error
}

type DBNoopAdapter interface {
	DB
	IsNoop() bool
}

type API struct {
	db         DBNoopAdapter
	url        *url.URL
	httpClient *resty.Client
}

func (a *API) SetDB(db DB) {
	a.db = DBNoopExtended{db: db}
}

func (a *API) MerkleRoots(ctx context.Context, opts ...queries.MerkleRootsQueryOption) (*queries.MerkleRootPage, error) {
	var query queries.MerkleRootsQuery
	for _, o := range opts {
		o(&query)
	}

	queryBuilder := querybuilders.NewQueryBuilder(querybuilders.WithFilterQueryBuilder(&merkleRootsFilterQueryBuilder{query: query}))
	params, err := queryBuilder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build merkle roots query params: %w", err)
	}

	var result queries.MerkleRootPage
	_, err = a.httpClient.R().
		SetContext(ctx).
		SetResult(&result).
		SetQueryParams(params.ParseToMap()).
		Get(a.url.String())
	if err != nil {
		return nil, fmt.Errorf("HTTP response failure: %w", err)
	}

	return &result, nil
}

func (a *API) SyncMerkleRoots(ctx context.Context) error {
	if a.db.IsNoop() {
		return nil
	}

	lastKey := a.db.GetLastMerkleRoot()
	prevKey := lastKey

	for {
		select {
		case <-ctx.Done():
			return customerr.ErrSyncMerkleRootsTimeout

		default:
			res, err := a.MerkleRoots(ctx, queries.MerkleRootsQueryWithLastEvaluatedKey(lastKey))
			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					return customerr.ErrSyncMerkleRootsTimeout
				}

				return fmt.Errorf("HTTP response failure: %w", err)
			}

			if len(res.Content) == 0 {
				return nil
			}

			lastKey = res.Page.LastEvaluatedKey
			if prevKey == lastKey {
				return customerr.ErrStaleLastEvaluatedKey
			}

			if err = a.db.SaveMerkleRoots(res.Content); err != nil {
				return fmt.Errorf("failed to save roots into db: %w", err)
			}

			if lastKey == "" {
				return nil
			}

			prevKey = lastKey
		}
	}
}

func NewAPI(baseURL *url.URL, httpClient *resty.Client) *API {
	return &API{
		url:        baseURL.JoinPath(route),
		httpClient: httpClient,
		db:         NoopDB{},
	}
}
