package main

import (
	"context"
	"log"

	wallet "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples/exampleutil"
)

func main() {
	spv, err := wallet.NewWithXPriv(exampleutil.ExampleConfig, exampleutil.ExampleXPriv)
	if err != nil {
		log.Fatal(err)
	}

	res, err := spv.SharedConfig(context.Background())
	if err != nil {
		log.Fatal()
	}
	exampleutil.Print("HTTP api/v1/configs/shared", res)
}
