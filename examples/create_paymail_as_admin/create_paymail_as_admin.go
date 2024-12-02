package main

import (
	"context"
	"log"

	wallet "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples/exampleutil"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/querybuilders"
	"github.com/bitcoin-sv/spv-wallet-go-client/walletkeys"
)

func main() {
	ctx := context.Background()
	keys, err := walletkeys.RandomKeys()
	if err != nil {
		log.Fatal(err)
	}

	adminAPI, err := wallet.NewAdminAPIWithXPriv(exampleutil.ExampleConfig, examples.XPriv)
	if err != nil {
		log.Fatal(err)
	}

	xPub, err := adminAPI.CreateXPub(ctx, &commands.CreateUserXpub{
		Metadata: map[string]any{"xpub_key": "xpub_val"},
		XPub:     keys.XPub(),
	})
	if err != nil {
		log.Fatal(err)
	}
	exampleutil.Print("[HTTP POST][Step 1] Create xPub - api/v1/admin/users", xPub)

	addr := exampleutil.RandomPaymail()
	paymail, err := adminAPI.CreatePaymail(ctx, &commands.CreatePaymail{
		Key:      keys.XPub(),
		Address:  addr,
		Metadata: querybuilders.Metadata{"key": "value"},
	})
	if err != nil {
		log.Fatal(err)
	}

	exampleutil.Print("[HTTP POST][Step 2] Create Paymail - api/v1/admin/paymails", paymail)
}
