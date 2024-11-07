package querybuilders_test

import (
	"net/url"
	"testing"
	"time"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/querybuilders"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/testfixtures"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/stretchr/testify/require"
)

func TestModelFilterQueryBuilder_Build(t *testing.T) {
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
				IncludeDeleted: testfixtures.Ptr(true),
			},
		},
		"model filter: filter with only created range 'from' field set": {
			expectedParams: url.Values{
				"createdRange[from]": []string{"2021-01-01T00:00:00Z"},
			},
			filter: filter.ModelFilter{
				CreatedRange: &filter.TimeRange{
					From: testfixtures.Ptr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
				},
			},
		},
		"model filter: filter wtth only created range 'to' field set": {
			expectedParams: url.Values{
				"createdRange[to]": []string{"2021-01-02T00:00:00Z"},
			},
			filter: filter.ModelFilter{
				CreatedRange: &filter.TimeRange{
					To: testfixtures.Ptr(time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)),
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
					From: testfixtures.Ptr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
					To:   testfixtures.Ptr(time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)),
				},
			},
		},
		"model filter: filter with only updated range 'from' field set": {
			expectedParams: url.Values{
				"updatedRange[from]": []string{"2021-02-01T00:00:00Z"},
			},
			filter: filter.ModelFilter{
				UpdatedRange: &filter.TimeRange{
					From: testfixtures.Ptr(time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)),
				},
			},
		},
		"model filter: filter with only updated range 'to' field set": {
			expectedParams: url.Values{
				"updatedRange[to]": []string{"2021-02-02T00:00:00Z"},
			},
			filter: filter.ModelFilter{
				UpdatedRange: &filter.TimeRange{
					To: testfixtures.Ptr(time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC)),
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
					From: testfixtures.Ptr(time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)),
					To:   testfixtures.Ptr(time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC)),
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
				IncludeDeleted: testfixtures.Ptr(true),
				CreatedRange: &filter.TimeRange{
					From: testfixtures.Ptr(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
					To:   testfixtures.Ptr(time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC)),
				},
				UpdatedRange: &filter.TimeRange{
					From: testfixtures.Ptr(time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)),
					To:   testfixtures.Ptr(time.Date(2021, 2, 2, 0, 0, 0, 0, time.UTC)),
				},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			m := querybuilders.ModelFilterBuilder{ModelFilter: tc.filter}
			got, err := m.Build()
			require.ErrorIs(t, tc.expectedErr, err)
			require.Equal(t, tc.expectedParams, got)
		})
	}
}
