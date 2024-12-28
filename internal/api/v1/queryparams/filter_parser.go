package queryparams

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"

	goclienterr "github.com/bitcoin-sv/spv-wallet-go-client/errors"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
)

type FilterParser struct {
	Filter any
}

func (f *FilterParser) Parse() (url.Values, error) {
	if f.Filter == nil {
		return nil, goclienterr.ErrNilFilterProvided
	}

	t := reflect.TypeOf(f.Filter)
	v := reflect.ValueOf(f.Filter)

	// Ensure Filter is a struct
	if t.Kind() != reflect.Struct {
		return nil, goclienterr.ErrFilterTypeNotStruct
	}

	params := NewURLValues()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// Extract the JSON tag
		tag := strings.Split(field.Tag.Get("json"), ",")[0]
		if tag == "" || value.IsZero() {
			continue
		}

		// Process the value based on its type
		switch field.Type {
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
			params.AddPair(tag, fmt.Sprintf("%v", value.Interface()))
		}
	}

	return params.Values, nil
}
