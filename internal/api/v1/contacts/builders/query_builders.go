package builders

import (
	"fmt"
	"net/url"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/querybuilders"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
)

type UserContactFilterQueryBuilder struct {
	ModelFilterBuilder querybuilders.ModelFilterBuilder
	ContactFilter      filter.ContactFilter
}

func (u *UserContactFilterQueryBuilder) Build() (url.Values, error) {
	modelFilterBuilder, err := u.ModelFilterBuilder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build model filter query params: %w", err)
	}

	params := querybuilders.NewExtendedURLValues()
	if len(modelFilterBuilder) > 0 {
		params.Append(modelFilterBuilder)
	}

	params.AddPair("id", u.ContactFilter.ID)
	params.AddPair("fullName", u.ContactFilter.FullName)
	params.AddPair("paymail", u.ContactFilter.Paymail)
	params.AddPair("pubKey", u.ContactFilter.PubKey)
	params.AddPair("status", u.ContactFilter.Status)
	return params.Values, nil
}
