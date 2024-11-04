package transactionsfixture

import (
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/testfixtures"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

func NewTransactions() *response.PageModel[response.Transaction] {
	return &response.PageModel[response.Transaction]{
		Content: []*response.Transaction{
			{
				Model: response.Model{
					CreatedAt: testfixtures.Parse("2024-10-07T14:03:26.736816Z"),
					UpdatedAt: testfixtures.Parse("2024-11-04T14:14:02.561175Z"),
					Metadata: map[string]any{
						"domain":          "john.doe.test.4chain.space",
						"example_key1":    "example_key10_val",
						"ip_address":      "127.0.0.01",
						"user_agent":      "node-fetch",
						"paymail_request": "HandleReceivedP2pTransaction",
						"reference_id":    "1c2dcc61-f48f-44f2-aba2-9a759a514d49",
						"p2p_tx_metadata": map[string]any{
							"pubkey": "3fa7af5b-4568-4873-86da-0aa442ca91dd",
							"sender": "john.doe@handcash.io",
						},
					},
				},
				ID:                   "2c250e21-c33a-41e3-a4e3-77c68b03244e",
				Hex:                  "283b1c6deb6d6263b3cec7a4701d46d3",
				XpubOutIDs:           []string{"4c9a0a0d-ea4f-4f03-b740-84438b3d210d"},
				BlockHash:            "47758f612c6bf5b454bcd642fe8031f6",
				BlockHeight:          512,
				Fee:                  1,
				NumberOfInputs:       2,
				NumberOfOutputs:      3,
				TotalValue:           311,
				OutputValue:          100,
				Status:               "MINED",
				TransactionDirection: "incoming",
			},
			{
				Model: response.Model{
					CreatedAt: testfixtures.Parse("2024-01-02T14:03:26.736816Z"),
					UpdatedAt: testfixtures.Parse("2024-01-04T14:14:02.561175Z"),
					Metadata: map[string]any{
						"domain":          "jane.doe.test.4chain.space",
						"example_key101":  "example_key101_val",
						"ip_address":      "127.0.0.01",
						"user_agent":      "node-fetch",
						"paymail_request": "HandleReceivedP2pTransaction",
						"reference_id":    "2c6dcc71-f42f-54f2-ada1-1c658a515d50",
						"p2p_tx_metadata": map[string]any{
							"pubkey": "4fa8af6b-3217-2373-76da-0aa552ca88aa",
							"sender": "jane.doe@handcash.io",
						},
					},
				},
				ID:                   "1c110e11-c23a-51e5-a7e7-99c12b01233e",
				Hex:                  "283b1c7deb7d7773b3cec7a8801d47d2",
				XpubOutIDs:           []string{"2c8a1a1d-ea5f-5f04-b890-92418b2d411d"},
				BlockHash:            "56659f622c6bf5b554bcd742fe8132f9",
				BlockHeight:          1024,
				Fee:                  1,
				NumberOfInputs:       2,
				NumberOfOutputs:      3,
				TotalValue:           500,
				OutputValue:          200,
				Status:               "MINED",
				TransactionDirection: "incoming",
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
