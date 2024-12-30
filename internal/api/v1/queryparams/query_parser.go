package queryparams

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"

	goclienterr "github.com/bitcoin-sv/spv-wallet-go-client/errors"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
)

type QueryParser[F queries.QueryFilters] struct {
	query *queries.Query[F]
}

func (q *QueryParser[F]) parse(v any) url.Values {
	params := NewURLValues()
	for i := 0; i < reflect.TypeOf(v).NumField(); i++ {
		field := reflect.TypeOf(v).Field(i)
		value := reflect.ValueOf(v).Field(i)

		// Extract the JSON tag
		tag := strings.Split(field.Tag.Get("json"), ",")[0]

		// Process the value based on its type
		switch field.Type {
		case reflect.TypeOf(filter.ModelFilter{}):
			params.Append(q.parse(value.Interface()))

		case reflect.PointerTo(reflect.TypeOf(filter.TimeRange{})):
			if !value.IsNil() {
				params.AddPair(tag, value.Interface().(*filter.TimeRange))
			}
		case reflect.PointerTo(reflect.TypeOf(false)): // Pointer to a bool
			if !value.IsNil() {
				params.AddPair(tag, value.Elem().Bool())
			}
		case reflect.PointerTo(reflect.TypeOf("")): // Pointer to a string
			if !value.IsNil() {
				params.AddPair(tag, value.Elem().String())
			}
		case reflect.PointerTo(reflect.TypeOf(uint64(0))), reflect.PointerTo(reflect.TypeOf(uint32(0))): // Pointer to a uint64, unit32
			if !value.IsNil() {
				params.AddPair(tag, value.Elem().Uint())
			}

		default:
			if value.IsValid() && value.CanInterface() {
				params.AddPair(tag, fmt.Sprintf("%v", value.Interface()))
			}
		}
	}

	return params.Values
}

func (q *QueryParser[F]) Parse() (*URLValues, error) {
	total := NewURLValues()
	if q.query.PageFilter != (filter.Page{}) {
		total.Append(q.parse(q.query.PageFilter))
	}

	metadata := &MetadataParser{Metadata: q.query.Metadata, MaxDepth: DefaultMaxDepth}
	params, err := metadata.Parse()
	if err != nil {
		return nil, fmt.Errorf("failed to parse metadata: %w", err)
	}
	total.Append(params)

	t := reflect.TypeOf(q.query.Filter)
	v := reflect.ValueOf(q.query.Filter)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// Extract the JSON tag
		tag := strings.Split(field.Tag.Get("json"), ",")[0]

		switch field.Type {
		case
			reflect.TypeOf(filter.UtxoFilter{}),
			reflect.TypeOf(filter.TransactionFilter{}),
			reflect.TypeOf(filter.PaymailFilter{}),
			reflect.TypeOf(filter.ContactFilter{}),
			reflect.TypeOf(filter.AccessKeyFilter{}),
			reflect.TypeOf(filter.XpubFilter{}):

			total.Append(q.parse(value.Interface()))

		case reflect.TypeOf(filter.ModelFilter{}):
			total.Append(q.parse(value.Interface()))

		case reflect.PointerTo(reflect.TypeOf(filter.TimeRange{})): // Pointer to a time range
			if !value.IsNil() {
				total.AddPair(tag, value.Interface().(*filter.TimeRange))
			}
		case reflect.PointerTo(reflect.TypeOf(false)): // Pointer to a bool
			if !value.IsNil() {
				total.AddPair(tag, value.Elem().Bool())
			}
		case reflect.PointerTo(reflect.TypeOf("")): // Pointer to a string
			if !value.IsNil() {
				total.AddPair(tag, value.Elem().String())
			}
		case reflect.PointerTo(reflect.TypeOf(uint64(0))), reflect.PointerTo(reflect.TypeOf(uint32(0))): // Pointer to a uint64, unit32
			if !value.IsNil() {
				total.AddPair(tag, value.Elem().Uint())
			}
		}
	}

	return total, nil
}

func NewQueryParser[F queries.QueryFilters](query *queries.Query[F]) (*QueryParser[F], error) {
	if query == nil {
		return nil, goclienterr.ErrQueryParserFailed
	}

	return &QueryParser[F]{query: query}, nil
}
