package testfixtures

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/spvwallettest"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

func ExpectedContactsPageFromUserAPI(t *testing.T) *queries.UserContactsPage {
	return &queries.UserContactsPage{
		Content: []*response.Contact{
			{
				Model: response.Model{
					CreatedAt: spvwallettest.ParseTime(t, "2024-10-18T12:07:44.739839Z"),
					UpdatedAt: spvwallettest.ParseTime(t, "2024-10-18T15:08:44.739918Z"),
				},
				ID:       "4f730efa-2a33-4275-bfdb-1f21fc110963",
				FullName: "John Doe",
				Paymail:  "john.doe.test5@john.doe.4chain.space",
				PubKey:   "19751ea9-6c1f-4ba7-a7e2-551ef7930136",
				Status:   "unconfirmed",
			},
			{
				Model: response.Model{
					CreatedAt: spvwallettest.ParseTime(t, "2024-10-18T12:07:44.739839Z"),
					UpdatedAt: spvwallettest.ParseTime(t, "2024-10-18T15:08:44.739918Z"),
				},
				ID:       "e55a4d4e-4a4b-4720-8556-1c00dd6a5cf3",
				FullName: "Jane Doe",
				Paymail:  "jane.doe.test5@jane.doe.4chain.space",
				PubKey:   "f8898969-3f96-48d3-b122-bbb3e738dbf5",
				Status:   "unconfirmed",
			},
		},
		Page: response.PageDescription{
			Size:          2,
			Number:        2,
			TotalElements: 2,
			TotalPages:    1,
		},
	}
}

func ExpectedContactWithPaymailFromUserAPI(t *testing.T) *response.Contact {
	return &response.Contact{
		Model: response.Model{
			CreatedAt: spvwallettest.ParseTime(t, "2024-10-18T12:07:44.739839Z"),
			UpdatedAt: spvwallettest.ParseTime(t, "2024-10-18T15:08:44.739918Z"),
		},
		ID:       "4f730efa-2a33-4275-bfdb-1f21fc110963",
		FullName: "John Doe",
		Paymail:  "john.doe.test5@john.doe.4chain.space",
		PubKey:   "19751ea9-6c1f-4ba7-a7e2-551ef7930136",
		Status:   "unconfirmed",
	}
}

func ExpectedUpsertContactFromUserAPI(t *testing.T) *response.Contact {
	return &response.Contact{
		Model: response.Model{
			CreatedAt: spvwallettest.ParseTime(t, "2024-10-18T12:07:44.739839Z"),
			UpdatedAt: spvwallettest.ParseTime(t, "2024-11-06T11:30:35.090124Z"),
			Metadata: map[string]interface{}{
				"example_key": "example_val",
			},
		},
		ID:       "68acf78f-5ece-4917-821d-8028ecf06c9a",
		FullName: "John Doe",
		Paymail:  "john.doe.test@john.doe.test.4chain.space",
		PubKey:   "0df36839-67bb-49e7-a9c7-e839aa564871",
		Status:   "unconfirmed",
	}
}

func ExpectedUpdatedUserContactFromAdminAPI(t *testing.T) *response.Contact {
	return &response.Contact{
		Model: response.Model{
			CreatedAt: spvwallettest.ParseTime(t, "2024-11-28T13:34:52.11722Z"),
			UpdatedAt: spvwallettest.ParseTime(t, "2024-11-29T08:23:19.66093Z"),
			Metadata:  map[string]any{"phoneNumber": "123456789"},
		},
		ID:       "4d570959-dd85-4f53-bad1-18d0671761e9",
		FullName: "John Doe Williams",
		Paymail:  "john.doe.test@john.doe.test.4chain.space",
		PubKey:   "96843af4-fc9c-4778-945d-2131ac5b1a8a",
		Status:   "awaiting",
	}
}

func ExpectedContactsPageFromAdminAPI(t *testing.T) *queries.UserContactsPage {
	return &queries.UserContactsPage{
		Content: []*response.Contact{
			{
				Model: response.Model{
					CreatedAt: spvwallettest.ParseTime(t, "2024-11-28T14:58:13.262238Z"),
					UpdatedAt: spvwallettest.ParseTime(t, "2024-11-28T16:18:43.842434Z"),
				},
				ID:       "7a5625ac-8256-454a-84a3-7f03f50cd7dc",
				FullName: "John Doe",
				Paymail:  "john.doe.test@john.doe.4chain.space",
				PubKey:   "bbbb7a4e-a3f4-4ca4-800a-fdd8029eda37",
				Status:   "confirmed",
			},
			{
				Model: response.Model{
					CreatedAt: spvwallettest.ParseTime(t, "2024-11-28T14:58:13.029966Z"),
					UpdatedAt: spvwallettest.ParseTime(t, "2024-11-28T14:58:13.03002Z"),
					Metadata: map[string]any{
						"phoneNumber": "123456789",
					},
				},
				ID:       "d05d2388-3c16-426d-98f1-ced9d9c5f4e1",
				FullName: "Jane Doe",
				Paymail:  "jane.doe.jane@john.doe.4chain.space",
				PubKey:   "ee191d63-1619-4fd3-ae3d-2202cfab751d",
				Status:   "unconfirmed",
			},
		},
		Page: response.PageDescription{
			Size:          50,
			Number:        1,
			TotalElements: 2,
			TotalPages:    1,
		},
	}
}
