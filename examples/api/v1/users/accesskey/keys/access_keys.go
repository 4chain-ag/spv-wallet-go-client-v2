package main

import (
	"context"
	"log"

	wallet "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples/exampleutil"
)

func main() {
	spv, err := wallet.NewWithXPriv(exampleutil.ExampleConfig, examples.XPriv)
	if err != nil {
		log.Fatal(err)
	}

	accessKeys, err := spv.AccessKeys(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	exampleutil.Print("[HTTP GET] Access keys - api/v1/users/current/keys/", accessKeys)
}
