package querybuilders

import (
	"net/url"

	"github.com/bitcoin-sv/spv-wallet/models/filter"
)

type TransactionFilterBuilder struct {
	TransactionFilter  filter.TransactionFilter
	ModelFilterBuilder ModelFilterBuilder
}

func (t *TransactionFilterBuilder) Build() (url.Values, error) {
	mfv, err := t.ModelFilterBuilder.Build()
	if err != nil {
		return nil, err
	}

	params := NewExtendedURLValues()
	if len(mfv) > 0 {
		params.Append(mfv)
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
