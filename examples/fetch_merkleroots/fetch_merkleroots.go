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

	page, err := spv.MerkleRoots(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	exampleutil.Print("[HTTP GET] Merkle roots page - api/v1/merkleroots", page)
}
