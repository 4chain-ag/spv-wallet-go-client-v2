package queryparams

import (
	"errors"
	"net/url"

	goclienterr "github.com/bitcoin-sv/spv-wallet-go-client/errors"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
)

// TODO: Handle Merkleroots search query params
// TODO: Handle Admin XPubs search query params
// TODO: Change types to be non-exported if not used outside pkg
// TODO: Update error docs

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
			return nil, errors.Join(err, goclienterr.ErrQueryParamsBuilder)
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
		parsers = append(parsers, &MetadataParser{
			Metadata: query.Metadata,
			MaxDepth: DefaultMaxDepth,
		})
	}

	var zero filter.Page
	if query.PageFilter != zero {
		parsers = append(parsers, &FilterParser{Filter: query.PageFilter})
	}

	switch f := any(query.Filter).(type) {
	case filter.AccessKeyFilter:
		parsers = append(parsers, &FilterParser{Filter: f.ModelFilter})
		parsers = append(parsers, &FilterParser{Filter: f})

	case filter.ContactFilter:
		parsers = append(parsers, &FilterParser{Filter: f.ModelFilter})
		parsers = append(parsers, &FilterParser{Filter: f})

	case filter.PaymailFilter:
		parsers = append(parsers, &FilterParser{Filter: f.ModelFilter})
		parsers = append(parsers, &FilterParser{Filter: f})

	case filter.TransactionFilter:
		parsers = append(parsers, &FilterParser{Filter: f.ModelFilter})
		parsers = append(parsers, &FilterParser{Filter: f})

	case filter.UtxoFilter:
		parsers = append(parsers, &FilterParser{Filter: f.ModelFilter})
		parsers = append(parsers, &FilterParser{Filter: f})

	case filter.AdminAccessKeyFilter:
		parsers = append(parsers, &FilterParser{Filter: f.AccessKeyFilter})
		parsers = append(parsers, &FilterParser{Filter: f.ModelFilter})
		parsers = append(parsers, &FilterParser{Filter: f})

	case filter.AdminContactFilter:
		parsers = append(parsers, &FilterParser{Filter: f.ContactFilter})
		parsers = append(parsers, &FilterParser{Filter: f.ModelFilter})
		parsers = append(parsers, &FilterParser{Filter: f})

	case filter.AdminPaymailFilter:
		parsers = append(parsers, &FilterParser{Filter: f.PaymailFilter})
		parsers = append(parsers, &FilterParser{Filter: f.ModelFilter})
		parsers = append(parsers, &FilterParser{Filter: f})

	case filter.AdminTransactionFilter:
		parsers = append(parsers, &FilterParser{Filter: f.TransactionFilter})
		parsers = append(parsers, &FilterParser{Filter: f.ModelFilter})
		parsers = append(parsers, &FilterParser{Filter: f})

	case filter.AdminUtxoFilter:
		parsers = append(parsers, &FilterParser{Filter: f.UtxoFilter})
		parsers = append(parsers, &FilterParser{Filter: f.ModelFilter})
		parsers = append(parsers, &FilterParser{Filter: f})
	}

	return &Builder{parsers: parsers}, nil
}
