package main

import (
	"context"
	"log"

	wallet "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/docs/examples"
	"github.com/bitcoin-sv/spv-wallet-go-client/docs/examples/exampleutil"
)

func main() {
	spv, err := wallet.NewWithXPriv(exampleutil.ExampleConfig, examples.XPriv)
	if err != nil {
		log.Fatal(err)
	}

	transaction, err := spv.RecordTransaction(context.Background(), &commands.RecordTransaction{
		Metadata:    map[string]any{"key": "value"},
		ReferenceID: "8bc53e34-b6fd-4e8b-b1b7-6f30f8f149f2",
		Hex:         "0100000002...",
	})
	if err != nil {
		log.Fatal(err)
	}

	exampleutil.Print("[HTTP POST] Record transaction - api/v1/transactions", transaction)
}
