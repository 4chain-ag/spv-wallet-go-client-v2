package transactions

import (
	"net/url"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/querybuilders"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
)

type TransactionFilterBuilder struct {
	TransactionFilter  filter.TransactionFilter
	ModelFilterBuilder querybuilders.ModelFilterBuilder
}

func (t *TransactionFilterBuilder) Build() (url.Values, error) {
	modelFilterBuilder, err := t.ModelFilterBuilder.Build()
	if err != nil {
		return nil, err
	}

	params := querybuilders.NewExtendedURLValues()
	if len(modelFilterBuilder) > 0 {
		params.Append(modelFilterBuilder)
	}

	params.AddPair("id", t.TransactionFilter.Id)
	params.AddPair("hex", t.TransactionFilter.Hex)
	params.AddPair("blockHash", t.TransactionFilter.BlockHash)
	params.AddPair("blockHeight", t.TransactionFilter.BlockHeight)
	params.AddPair("fee", t.TransactionFilter.Fee)
	params.AddPair("numberOfInputs", t.TransactionFilter.NumberOfInputs)
	params.AddPair("numberOfOutputs", t.TransactionFilter.NumberOfOutputs)
	params.AddPair("draftId", t.TransactionFilter.DraftID)
	params.AddPair("totalValue", t.TransactionFilter.TotalValue)
	params.AddPair("status", t.TransactionFilter.Status)
	return params.Values, nil
}