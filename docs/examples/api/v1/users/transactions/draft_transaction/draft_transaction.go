package main

import (
	"context"
	"log"

	wallet "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/docs/examples"
	"github.com/bitcoin-sv/spv-wallet-go-client/docs/examples/exampleutil"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

func main() {
	spv, err := wallet.NewWithXPriv(exampleutil.ExampleConfig, examples.XPriv)
	if err != nil {
		log.Fatal(err)
	}

	transaction, err := spv.DraftTransaction(context.Background(), &commands.DraftTransaction{
		Config: response.TransactionConfig{
			Outputs: []*response.TransactionOutput{
				{
					To:       "receiver@example.com",
					Satoshis: 1,
				},
			},
		},
		Metadata: map[string]any{"key": "value"},
	})
	if err != nil {
		log.Fatal(err)
	}

	exampleutil.Print("[HTTP POST] Draft transaction - api/v1/transactions", transaction)
}
