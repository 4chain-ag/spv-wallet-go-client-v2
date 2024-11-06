package querybuilders

import (
	"errors"
	"net/url"

	"github.com/bitcoin-sv/spv-wallet/models/filter"
)

type QueryBuilderOption func(*QueryBuilder)

func WithQueryParamsFilter(q filter.QueryParams) QueryBuilderOption {
	var zero filter.QueryParams
	return func(qb *QueryBuilder) {
		if q != zero {
			qb.builders = append(qb.builders, &QueryParamsFilterQueryBuilder{q})
		}
	}
}

func WithMetadataFilter(m Metadata) QueryBuilderOption {
	return func(qb *QueryBuilder) {
		if m != nil {
			qb.builders = append(qb.builders, &MetadataFilterQueryBuilder{MaxDepth: DefaultMaxDepth, Metadata: m})
		}
	}
}

func WithTransactionFilter(tf filter.TransactionFilter) QueryBuilderOption {
	var zero filter.TransactionFilter
	return func(qb *QueryBuilder) {
		if tf != zero {
			qb.builders = append(qb.builders, &TransactionFilterQueryBuilder{
				TransactionFilter:       tf,
				ModelFilterQueryBuilder: ModelFilterQueryBuilder{ModelFilter: tf.ModelFilter},
			})
		}
	}
}

func WithFilterQueryBuilder(b FilterQueryBuilder) QueryBuilderOption {
	return func(qb *QueryBuilder) {
		if b != nil {
			qb.builders = append(qb.builders, b)
		}
	}
}

type FilterQueryBuilder interface {
	Build() (url.Values, error)
}

type QueryBuilder struct {
	builders []FilterQueryBuilder
}

func (q *QueryBuilder) Build() (url.Values, error) {
	params := NewExtendedURLValues()
	for _, b := range q.builders {
		bparams, err := b.Build()
		if err != nil {
			return nil, errors.Join(err, ErrFilterQueryBuilder)
		}

		if len(bparams) > 0 {
			params.Append(bparams)
		}
	}
	return params.Values, nil
}

func NewQueryBuilder(opts ...QueryBuilderOption) *QueryBuilder {
	var qb QueryBuilder
	for _, o := range opts {
		o(&qb)
	}
	return &qb
}

var ErrFilterQueryBuilder = errors.New("filter query builder - build operation failure")

func ParseToMap(values url.Values) map[string]string {
	m := make(map[string]string)
	for k, v := range values {
		m[k] = v[0]
	}
	return m
}
