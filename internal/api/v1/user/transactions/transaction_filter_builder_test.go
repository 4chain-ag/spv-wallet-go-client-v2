package transactions_test

import (
	"net/url"
	"testing"
	"time"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/querybuilders"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/transactions"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/testfixtures"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/stretchr/testify/require"
)

func TestTransactionFilterBuilder_Build(t *testing.T) {
	tests := map[string]struct {
		filter         filter.TransactionFilter
		expectedParams url.Values
		expectedErr    error
	}{
		"transaction filter: zero values": {
			filter: filter.TransactionFilter{
				Id:              testfixtures.Ptr(""),
				Hex:             testfixtures.Ptr(""),
				BlockHash:       testfixtures.Ptr(""),
				BlockHeight:     testfixtures.Ptr(uint64(0)),
				Fee:             testfixtures.Ptr(uint64(0)),
				NumberOfInputs:  testfixtures.Ptr(uint32(0)),
				NumberOfOutputs: testfixtures.Ptr(uint32(0)),
				DraftID:         testfixtures.Ptr(""),
				TotalValue:      testfixtures.Ptr(uint64(0)),
				Status:          testfixtures.Ptr(""),
			},
			expectedParams: make(url.Values),
		},
		"transaction filter: filter with only 'id' field set": {
			filter: filter.TransactionFilter{
				Id: testfixtures.Ptr("d425432e0d10a46af1ec6d00f380e9581ebf7907f3486572b3cd561a4c326e14"),
			},
			expectedParams: url.Values{
				"id": []string{"d425432e0d10a46af1ec6d00f380e9581ebf7907f3486572b3cd561a4c326e14"},
			},
		},
		"transaction filter: filter with only 'hex' field set": {
			filter: filter.TransactionFilter{
				Hex: testfixtures.Ptr("001290b87619e679aaf6b8aadd30c778726c89fc4442110feb6d8265a190386beb8311a31e7e97a1c9ff2c84f3993283078965eb81f6fa64f3d7ba7fdd09678d"),
			},
			expectedParams: url.Values{
				"hex": []string{"001290b87619e679aaf6b8aadd30c778726c89fc4442110feb6d8265a190386beb8311a31e7e97a1c9ff2c84f3993283078965eb81f6fa64f3d7ba7fdd09678d"},
			},
		},
		"transaction filter: filter with only 'block hash' field set": {
			filter: filter.TransactionFilter{
				BlockHash: testfixtures.Ptr("0000000000000000031928c28075a82d7a00c2c90b489d1d66dc0afa3f8d26f8"),
			},
			expectedParams: url.Values{
				"blockHash": []string{"0000000000000000031928c28075a82d7a00c2c90b489d1d66dc0afa3f8d26f8"},
			},
		},
		"transaction filter: filter with only 'block height' field set": {
			filter: filter.TransactionFilter{
				BlockHeight: testfixtures.Ptr(uint64(839376)),
			},
			expectedParams: url.Values{
				"blockHeight": []string{"839376"},
			},
		},
		"transaction filter: filter with only 'fee' field set": {
			filter: filter.TransactionFilter{
				Fee: testfixtures.Ptr(uint64(1)),
			},
			expectedParams: url.Values{
				"fee": []string{"1"},
			},
		},
		"transaction filter: filter with only 'number of inputs' field set": {
			filter: filter.TransactionFilter{
				NumberOfInputs: testfixtures.Ptr(uint32(10)),
			},
			expectedParams: url.Values{
				"numberOfInputs": []string{"10"},
			},
		},
		"transaction filter: filter with only 'number of outputs' field set": {
			filter: filter.TransactionFilter{
				NumberOfOutputs: testfixtures.Ptr(uint32(20)),
			},
			expectedParams: url.Values{
				"numberOfOutputs": []string{"20"},
			},
		},
		"transaction filter: filter with only 'draft id' field set": {
			filter: filter.TransactionFilter{
				DraftID: testfixtures.Ptr("d425432e0d10a46af1ec6d00f380e9581ebf7907f3486572b3cd561a4c326e14"),
			},
			expectedParams: url.Values{
				"draftId": []string{"d425432e0d10a46af1ec6d00f380e9581ebf7907f3486572b3cd561a4c326e14"},
			},
		},
		"transaction filter: filter with only 'total value' field set": {
			filter: filter.TransactionFilter{
				TotalValue: testfixtures.Ptr(uint64(100000000)),
			},
			expectedParams: url.Values{
				"totalValue": []string{"100000000"},
			},
		},
		"transaction filter: filter with only 'status' field set": {
			filter: filter.TransactionFilter{
				Status: testfixtures.Ptr("RECEIVED"),
			},
			expectedParams: url.Values{
				"status": []string{"RECEIVED"},
			},
		},
		"transaction filter: filter with only 'model filter' fields set": {
			filter: filter.TransactionFilter{
				ModelFilter: filter.ModelFilter{
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
			expectedParams: url.Values{
				"includeDeleted":     []string{"true"},
				"createdRange[from]": []string{"2021-01-01T00:00:00Z"},
				"createdRange[to]":   []string{"2021-01-02T00:00:00Z"},
				"updatedRange[from]": []string{"2021-02-01T00:00:00Z"},
				"updatedRange[to]":   []string{"2021-02-02T00:00:00Z"},
			},
		},
		"transaction filter: all fields set": {
			filter: filter.TransactionFilter{
				Id:              testfixtures.Ptr("d425432e0d10a46af1ec6d00f380e9581ebf7907f3486572b3cd561a4c326e14"),
				Hex:             testfixtures.Ptr("001290b87619e679aaf6b8aadd30c778726c89fc4442110feb6d8265a190386beb8311a31e7e97a1c9ff2c84f3993283078965eb81f6fa64f3d7ba7fdd09678d"),
				BlockHash:       testfixtures.Ptr("0000000000000000031928c28075a82d7a00c2c90b489d1d66dc0afa3f8d26f8"),
				BlockHeight:     testfixtures.Ptr(uint64(839376)),
				Fee:             testfixtures.Ptr(uint64(1)),
				NumberOfInputs:  testfixtures.Ptr(uint32(10)),
				NumberOfOutputs: testfixtures.Ptr(uint32(20)),
				DraftID:         testfixtures.Ptr("d425432e0d10a46af1ec6d00f380e9581ebf7907f3486572b3cd561a4c326e14"),
				TotalValue:      testfixtures.Ptr(uint64(100000000)),
				Status:          testfixtures.Ptr("RECEIVED"),
				ModelFilter: filter.ModelFilter{
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
			expectedParams: url.Values{
				"id":                 []string{"d425432e0d10a46af1ec6d00f380e9581ebf7907f3486572b3cd561a4c326e14"},
				"hex":                []string{"001290b87619e679aaf6b8aadd30c778726c89fc4442110feb6d8265a190386beb8311a31e7e97a1c9ff2c84f3993283078965eb81f6fa64f3d7ba7fdd09678d"},
				"blockHash":          []string{"0000000000000000031928c28075a82d7a00c2c90b489d1d66dc0afa3f8d26f8"},
				"blockHeight":        []string{"839376"},
				"fee":                []string{"1"},
				"numberOfInputs":     []string{"10"},
				"numberOfOutputs":    []string{"20"},
				"draftId":            []string{"d425432e0d10a46af1ec6d00f380e9581ebf7907f3486572b3cd561a4c326e14"},
				"totalValue":         []string{"100000000"},
				"status":             []string{"RECEIVED"},
				"includeDeleted":     []string{"true"},
				"createdRange[from]": []string{"2021-01-01T00:00:00Z"},
				"createdRange[to]":   []string{"2021-01-02T00:00:00Z"},
				"updatedRange[from]": []string{"2021-02-01T00:00:00Z"},
				"updatedRange[to]":   []string{"2021-02-02T00:00:00Z"},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tfb := transactions.TransactionFilterBuilder{
				TransactionFilter:  tc.filter,
				ModelFilterBuilder: querybuilders.ModelFilterBuilder{ModelFilter: tc.filter.ModelFilter},
			}
			got, err := tfb.Build()
			require.ErrorIs(t, tc.expectedErr, err)
			require.Equal(t, tc.expectedParams, got)
		})
	}
}