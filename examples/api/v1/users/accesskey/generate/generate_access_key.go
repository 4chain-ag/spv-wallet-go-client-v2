package main

import (
	"context"
	"log"

	wallet "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples/exampleutil"
)

func main() {
	spv, err := wallet.NewWithXPriv(exampleutil.ExampleConfig, examples.XPriv)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	cmd := commands.GenerateAccessKey{Metadata: map[string]any{"example_key": "example_value"}}
	key, err := spv.GenerateAccessKey(ctx, &cmd)
	if err != nil {
		log.Fatal(err)
	}

	exampleutil.Print("[HTTP POST] Generate access key - api/v1/users/current/keys", key)
}
