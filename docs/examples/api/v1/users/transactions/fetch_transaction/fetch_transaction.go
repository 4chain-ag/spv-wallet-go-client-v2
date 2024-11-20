package main

import (
	"context"
	"fmt"
	"log"

	wallet "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/docs/examples"
	"github.com/bitcoin-sv/spv-wallet-go-client/docs/examples/exampleutil"
)

func main() {
	spv, err := wallet.NewWithXPriv(exampleutil.ExampleConfig, examples.XPriv)
	if err != nil {
		log.Fatal(err)
	}

	transactionID := "2a4fa5b8caa54e0e46e8147b389e8efa09eb453ad8bb2577c56d67032a985e74"
	transaction, err := spv.Transaction(context.Background(), transactionID)
	if err != nil {
		log.Fatal(err)
	}

	exampleutil.Print(fmt.Sprintf("[HTTP GET] Transaction - api/v1/transactions/%s", transactionID), transaction)
}
