package builders

import (
	"fmt"
	"net/url"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/querybuilders"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
)

type AdminAccessKeyFilterQueryBuilder struct {
	AccessKeyFilter filter.AdminAccessKeyFilter
}

func (a *AdminAccessKeyFilterQueryBuilder) Build() (url.Values, error) {
	builder := UserAccessKeyFilterQueryBuilder{
		AccessKeyFilter:    a.AccessKeyFilter.AccessKeyFilter,
		ModelFilterBuilder: querybuilders.ModelFilterBuilder{ModelFilter: a.AccessKeyFilter.ModelFilter},
	}
	params, err := builder.BuildExtendedURLValues()
	if err != nil {
		return nil, fmt.Errorf("failed to build extended URL values: %w", err)
	}

	params.AddPair("xpubId", a.AccessKeyFilter.XpubID)
	return params.Values, nil
}

type UserAccessKeyFilterQueryBuilder struct {
	AccessKeyFilter    filter.AccessKeyFilter
	ModelFilterBuilder querybuilders.ModelFilterBuilder
}

func (u *UserAccessKeyFilterQueryBuilder) BuildExtendedURLValues() (*querybuilders.ExtendedURLValues, error) {
	modelFilterBuilder, err := u.ModelFilterBuilder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build model filter query params: %w", err)
	}

	params := querybuilders.NewExtendedURLValues()
	if len(modelFilterBuilder) > 0 {
		params.Append(modelFilterBuilder)
	}

	params.AddPair("revokedRange", u.AccessKeyFilter.RevokedRange)
	return params, nil
}

func (u *UserAccessKeyFilterQueryBuilder) Build() (url.Values, error) {
	params, err := u.BuildExtendedURLValues()
	if err != nil {
		return nil, fmt.Errorf("failed to build extended URL values: %w", err)
	}

	return params.Values, nil
}
