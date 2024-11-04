package query

import (
	"errors"
	"net/url"

	"github.com/bitcoin-sv/spv-wallet/models/filter"
)

type BuilderOption func(*Builder)

func WithQueryParamsFilter(q filter.QueryParams) BuilderOption {
	return func(b *Builder) {
		b.builders = append(b.builders, &QueryParamsFilterQueryBuilder{q})
	}
}

func WithMetadataFilter(m Metadata) BuilderOption {
	return func(b *Builder) {
		b.builders = append(b.builders, &MetadataFilterQueryBuilder{MaxDepth: DefaultMaxDepth, Metadata: m})
	}
}

func WithTransactionFilter(tf filter.TransactionFilter) BuilderOption {
	return func(b *Builder) {
		b.builders = append(b.builders, &TransactionFilterQueryBuilder{
			TransactionFilter:       tf,
			ModelFilterQueryBuilder: ModelFilterQueryBuilder{ModelFilter: tf.ModelFilter},
		})
	}
}

func WithFilterQueryBuilder(b FilterQueryBuilder) BuilderOption {
	return func(b *Builder) {
		if b != nil {
			b.builders = append(b.builders, b)
		}
	}
}

type FilterQueryBuilder interface {
	Build() (url.Values, error)
}

type Builder struct {
	builders []FilterQueryBuilder
}

func (q *Builder) Build() (url.Values, error) {
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

func NewQueryBuilder(opts ...BuilderOption) *Builder {
	var qb Builder
	for _, o := range opts {
		o(&qb)
	}
	return &qb
}

func Parse(values url.Values) map[string]string {
	m := make(map[string]string)
	for k, v := range values {
		m[k] = v[0]
	}
	return m
}

var ErrFilterQueryBuilder = errors.New("filter query builder - build query failure")
