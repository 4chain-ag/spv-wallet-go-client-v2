package contacts_test

import (
	"net/url"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/querybuilders"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/contacts"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/testutils"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/stretchr/testify/require"
)

func TestContactFilterQueryBuilder_Build(t *testing.T) {
	tests := map[string]struct {
		filter         filter.ContactFilter
		expectedParams url.Values
		expectedErr    error
	}{
		"contact filter: zero values": {
			expectedParams: make(url.Values),
		},
		"contact filter: filter with only 'id' field set": {
			filter: filter.ContactFilter{
				ID: testutils.Ptr("e3a1e174-cdf8-4683-b112-e198144eb9d2"),
			},
			expectedParams: url.Values{
				"id": []string{"e3a1e174-cdf8-4683-b112-e198144eb9d2"},
			},
		},
		"contact filter: filter with only 'full name' field set": {
			filter: filter.ContactFilter{
				FullName: testutils.Ptr("John Doe"),
			},
			expectedParams: url.Values{
				"fullName": []string{"John Doe"},
			},
		},
		"contact filter: filter with only 'paymail' field set": {
			filter: filter.ContactFilter{
				Paymail: testutils.Ptr("john.doe@test.com"),
			},
			expectedParams: url.Values{
				"paymail": []string{"john.doe@test.com"},
			},
		},
		"contact filter: filter with only 'status' field set": {
			filter: filter.ContactFilter{
				Status: testutils.Ptr("confirmed"),
			},
			expectedParams: url.Values{
				"status": []string{"confirmed"},
			},
		},
		"contact filter: filter with all fields set": {
			filter: filter.ContactFilter{
				ID:       testutils.Ptr("e3a1e174-cdf8-4683-b112-e198144eb9d2"),
				FullName: testutils.Ptr("John Doe"),
				Paymail:  testutils.Ptr("john.doe@test.com"),
				Status:   testutils.Ptr("confirmed"),
			},
			expectedParams: url.Values{
				"paymail":  []string{"john.doe@test.com"},
				"status":   []string{"confirmed"},
				"id":       []string{"e3a1e174-cdf8-4683-b112-e198144eb9d2"},
				"fullName": []string{"John Doe"},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			queryBuilder := contacts.ContactFilterQueryBuilder{
				ContactFilter:      tc.filter,
				ModelFilterBuilder: querybuilders.ModelFilterBuilder{ModelFilter: tc.filter.ModelFilter},
			}

			// then:
			got, err := queryBuilder.Build()
			require.ErrorIs(t, tc.expectedErr, err)
			require.Equal(t, tc.expectedParams, got)
		})
	}
}
