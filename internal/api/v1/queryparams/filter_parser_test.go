package queryparams_test

import (
	"net/url"
	"testing"
	"time"

	goclienterr "github.com/bitcoin-sv/spv-wallet-go-client/errors"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/queryparams"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/testutils"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/stretchr/testify/require"
)

func TestFilterParser_Parse(t *testing.T) {
	tests := map[string]struct {
		filter         any
		expectedParams url.Values
		expectedErr    error
	}{
		"filter: empty struct": {
			filter:         struct{}{},
			expectedParams: make(url.Values),
		},
		"filter: non-struct input": {
			filter:      "abcd",
			expectedErr: goclienterr.ErrFilterTypeNotStruct,
		},
		"filter: nil input": {
			expectedErr: goclienterr.ErrNilFilterProvided,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			parser := queryparams.FilterParser{Filter: tc.filter}

			// when:
			got, err := parser.Parse()

			// then:
			require.ErrorIs(t, tc.expectedErr, err)
			require.Equal(t, tc.expectedParams, got)
		})
	}
}

func TestFilterParser_Parse_ModelFilter(t *testing.T) {
	tests := map[string]struct {
		filter         filter.ModelFilter
		expectedParams url.Values
		expectedErr    error
	}{
		"model filter: filter with only 'include deleted field set": {
			expectedParams: url.Values{
				"includeDeleted": []string{"true"},
			},
			filter: filter.ModelFilter{
				IncludeDeleted: testutils.Ptr(true),
			},
		},
		"model filter: filter with only created range 'from' field set": {
			expectedParams: url.Values{
				"createdRange[from]": []string{"2021-01-01T00:00:00Z"},
			},
			filter: filter.ModelFilter{
				CreatedRange: &filter.TimeRange{
					From: testutils.Ptr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
			},
		},
		"model filter: filter wtth only created range 'to' field set": {
			expectedParams: url.Values{
				"createdRange[to]": []string{"2021-01-02T00:00:00Z"},
			},
			filter: filter.ModelFilter{
				CreatedRange: &filter.TimeRange{
					To: testutils.Ptr(time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)),
				},
			},
		},
		"model filter: filter with only created range both fields set": {
			expectedParams: url.Values{
				"createdRange[from]": []string{"2021-01-01T00:00:00Z"},
				"createdRange[to]":   []string{"2021-01-02T00:00:00Z"},
			},
			filter: filter.ModelFilter{
				CreatedRange: &filter.TimeRange{
					From: testutils.Ptr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
					To:   testutils.Ptr(time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)),
				},
			},
		},
		"model filter: filter with only updated range 'from' field set": {
			expectedParams: url.Values{
				"updatedRange[from]": []string{"2021-02-01T00:00:00Z"},
			},
			filter: filter.ModelFilter{
				UpdatedRange: &filter.TimeRange{
					From: testutils.Ptr(time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)),
				},
			},
		},
		"model filter: filter with only updated range 'to' field set": {
			expectedParams: url.Values{
				"updatedRange[to]": []string{"2021-02-02T00:00:00Z"},
			},
			filter: filter.ModelFilter{
				UpdatedRange: &filter.TimeRange{
					To: testutils.Ptr(time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC)),
				},
			},
		},
		"model filter: filter with only updated range both fields set": {
			expectedParams: url.Values{
				"updatedRange[from]": []string{"2021-02-01T00:00:00Z"},
				"updatedRange[to]":   []string{"2021-02-02T00:00:00Z"},
			},
			filter: filter.ModelFilter{
				UpdatedRange: &filter.TimeRange{
					From: testutils.Ptr(time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)),
					To:   testutils.Ptr(time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC)),
				},
			},
		},
		"model filter: all fields set": {
			expectedParams: url.Values{
				"includeDeleted":     []string{"true"},
				"createdRange[from]": []string{"2021-01-01T00:00:00Z"},
				"createdRange[to]":   []string{"2021-01-02T00:00:00Z"},
				"updatedRange[from]": []string{"2021-02-01T00:00:00Z"},
				"updatedRange[to]":   []string{"2021-02-02T00:00:00Z"},
			},
			filter: filter.ModelFilter{
				IncludeDeleted: testutils.Ptr(true),
				CreatedRange: &filter.TimeRange{
					From: testutils.Ptr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
					To:   testutils.Ptr(time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)),
				},
				UpdatedRange: &filter.TimeRange{
					From: testutils.Ptr(time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)),
					To:   testutils.Ptr(time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC)),
				},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			parser := queryparams.FilterParser{Filter: tc.filter}

			// when:
			got, err := parser.Parse()

			// then:
			require.ErrorIs(t, tc.expectedErr, err)
			require.Equal(t, tc.expectedParams, got)
		})
	}
}

func TestFilterParser_Parse_PageFilter(t *testing.T) {
	tests := map[string]struct {
		filter         filter.Page
		expectedParams url.Values
		expectedErr    error
	}{
		"page filter: filter with only 'number' set": {
			filter: filter.Page{
				Number: 10,
			},
			expectedParams: url.Values{
				"page": []string{"10"},
			},
		},
		"page filter: filter with only 'size' set": {
			filter: filter.Page{
				Size: 20,
			},
			expectedParams: url.Values{
				"size": []string{"20"},
			},
		},
		"page filter: filter with only 'sort' set": {
			filter: filter.Page{
				Sort: "asc",
			},
			expectedParams: url.Values{
				"sort": []string{"asc"},
			},
		},
		"page filter: filter with only 'sortBy' set": {
			filter: filter.Page{
				SortBy: "key",
			},
			expectedParams: url.Values{
				"sortBy": []string{"key"},
			},
		},
		"page filter: all fields set": {
			filter: filter.Page{
				Number: 10,
				Size:   20,
				Sort:   "asc",
				SortBy: "key",
			},
			expectedParams: url.Values{
				"sortBy": []string{"key"},
				"sort":   []string{"asc"},
				"size":   []string{"20"},
				"page":   []string{"10"},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			parser := queryparams.FilterParser{Filter: tc.filter}

			// when:
			got, err := parser.Parse()

			// then:
			require.ErrorIs(t, tc.expectedErr, err)
			require.Equal(t, tc.expectedParams, got)
		})
	}
}
