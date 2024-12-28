package queryparams

import (
	"errors"
	"net/url"

	goclienterr "github.com/bitcoin-sv/spv-wallet-go-client/errors"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
)

type Parser interface {
	Parse() (url.Values, error)
}

type Builder struct {
	parsers []Parser
}

func (b *Builder) Build() (*URLValues, error) {
	params := NewURLValues()
	for _, p := range b.parsers {
		values, err := p.Parse()
		if err != nil {
			return nil, errors.Join(err, goclienterr.ErrQueryParamsBuildFailed)
		}

		if len(values) > 0 {
			params.Append(values)
		}
	}

	return params, nil
}

func NewBuilder[F queries.QueryFilters](query *queries.Query[F]) (*Builder, error) {
	var parsers []Parser
	if query.Metadata != nil {
		parsers = append(parsers, &MetadataParser{Metadata: query.Metadata, MaxDepth: DefaultMaxDepth})
	}

	if query.PageFilter != (filter.Page{}) {
		parsers = append(parsers, &FilterParser{Filter: query.PageFilter})
	}

	parsers = append(parsers, initNonAdminQueryFilters(query)...)
	parsers = append(parsers, initAdminQueryFilters(query)...)
	return &Builder{parsers: parsers}, nil
}

func initNonAdminQueryFilters[F queries.QueryFilters](query *queries.Query[F]) []Parser {
	switch f := any(query.Filter).(type) {
	case filter.XpubFilter:
		return []Parser{&FilterParser{Filter: f.ModelFilter}, &FilterParser{Filter: f}}

	case filter.AccessKeyFilter:
		return []Parser{&FilterParser{Filter: f.ModelFilter}, &FilterParser{Filter: f}}

	case filter.ContactFilter:
		return []Parser{&FilterParser{Filter: f.ModelFilter}, &FilterParser{Filter: f}}

	case filter.PaymailFilter:
		return []Parser{&FilterParser{Filter: f.ModelFilter}, &FilterParser{Filter: f}}

	case filter.TransactionFilter:
		return []Parser{&FilterParser{Filter: f.ModelFilter}, &FilterParser{Filter: f}}

	case filter.UtxoFilter:
		return []Parser{&FilterParser{Filter: f.ModelFilter}, &FilterParser{Filter: f}}
	}

	return nil
}

func initAdminQueryFilters[F queries.QueryFilters](query *queries.Query[F]) []Parser {
	switch f := any(query.Filter).(type) {
	case filter.AdminAccessKeyFilter:
		return []Parser{
			&FilterParser{Filter: f.AccessKeyFilter},
			&FilterParser{Filter: f.ModelFilter},
			&FilterParser{Filter: f}}

	case filter.AdminContactFilter:
		return []Parser{
			&FilterParser{Filter: f.ContactFilter},
			&FilterParser{Filter: f.ModelFilter},
			&FilterParser{Filter: f}}

	case filter.AdminPaymailFilter:
		return []Parser{
			&FilterParser{Filter: f.PaymailFilter},
			&FilterParser{Filter: f.ModelFilter},
			&FilterParser{Filter: f}}

	case filter.AdminTransactionFilter:
		return []Parser{
			&FilterParser{Filter: f.TransactionFilter},
			&FilterParser{Filter: f.ModelFilter},
			&FilterParser{Filter: f}}

	case filter.AdminUtxoFilter:
		return []Parser{
			&FilterParser{Filter: f.UtxoFilter},
			&FilterParser{Filter: f.ModelFilter},
			&FilterParser{Filter: f}}
	}

	return nil
}
