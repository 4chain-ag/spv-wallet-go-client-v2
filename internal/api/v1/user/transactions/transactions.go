package transactions

import (
	"context"
	"fmt"

	"github.com/bitcoin-sv/spv-wallet-go-client/query"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/go-resty/resty/v2"
)

const route = "api/v1/transactions"

type DraftTransactionRequest struct {
	Config   response.TransactionConfig `json:"config"`
	Metadata query.Metadata             `json:"metadata"`
}

type RecordTransactionRequest struct {
	Metadata    query.Metadata `json:"metadata"`
	Hex         string         `json:"hex"`
	ReferenceID string         `json:"referenceId"`
}

type UpdateTransactionMetadataRequest struct {
	ID       string         `json:"-"`
	Metadata query.Metadata `json:"metadata"`
}

type API struct {
	addr string
	cli  *resty.Client
}

func (a *API) DraftTransaction(ctx context.Context, r DraftTransactionRequest) (*response.DraftTransaction, error) {
	var result response.DraftTransaction

	URL := a.addr + "/drafts"
	_, err := a.cli.R().
		SetContext(ctx).
		SetResult(&result).
		SetBody(r).
		Post(URL)
	if err != nil {
		return nil, fmt.Errorf("HTTP response failure: %w", err)
	}

	return &result, nil
}

func (a *API) RecordTransaction(ctx context.Context, r RecordTransactionRequest) (*response.Transaction, error) {
	var result response.Transaction

	_, err := a.cli.R().
		SetContext(ctx).
		SetResult(&result).
		SetBody(r).
		Patch(a.addr)
	if err != nil {
		return nil, fmt.Errorf("HTTP response failure: %w", err)
	}

	return &result, nil
}

func (a *API) UpdateTransactionMetadata(ctx context.Context, r UpdateTransactionMetadataRequest) (*response.Transaction, error) {
	var result response.Transaction

	URL := a.addr + "/" + r.ID
	_, err := a.cli.R().
		SetContext(ctx).
		SetResult(&result).
		SetBody(r).
		Patch(URL)
	if err != nil {
		return nil, fmt.Errorf("HTTP response failure: %w", err)
	}

	return &result, nil
}

func (a *API) Transaction(ctx context.Context, ID string) (*response.Transaction, error) {
	var result response.Transaction

	URL := a.addr + "/" + ID
	_, err := a.cli.R().
		SetContext(ctx).
		SetResult(&result).
		Get(URL)
	if err != nil {
		return nil, fmt.Errorf("HTTP response failure: %w", err)
	}

	return &result, nil
}

func (a *API) Transactions(ctx context.Context, opts ...query.BuilderOption) ([]*response.Transaction, error) {
	params, err := query.NewQueryBuilder(opts...).Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create transactions query params: %w", err)
	}

	var result response.PageModel[response.Transaction]
	_, err = a.cli.
		R().
		SetContext(ctx).
		SetResult(&result).
		SetQueryParams(query.Parse(params)).
		Get(a.addr)
	if err != nil {
		return nil, fmt.Errorf("HTTP response failure: %w", err)
	}

	return result.Content, nil
}

func NewAPI(addr string, cli *resty.Client) *API {
	return &API{
		addr: addr + "/" + route,
		cli:  cli,
	}
}
