package main

import (
	"context"
	"fmt"
	"log"

	wallet "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples/exampleutil"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/queryparams"
	"github.com/bitcoin-sv/spv-wallet-go-client/walletkeys"
)

func main() {
	keys, err := walletkeys.RandomKeys()
	if err != nil {
		log.Fatalf("Failed to generate random keys: %v", err)
	}
	fmt.Printf("Generated xPub for user: %s\n", keys.XPub())

	adminAPI, err := wallet.NewAdminAPIWithXPub(exampleutil.NewDefaultConfig(), examples.AdminXPriv)
	if err != nil {
		log.Fatalf("Failed to initialize admin API with XPriv: %v", err)
	}

	ctx := context.Background()
	xPub, err := adminAPI.CreateXPub(ctx, &commands.CreateUserXpub{
		XPub:     keys.XPub(),
		Metadata: queryparams.Metadata{"key": "value"},
	})
	if err != nil {
		log.Fatalf("Failed to create xPub: %v", err)
	}
	exampleutil.PrettyPrint("Created XPub", xPub)

	paymail, err := adminAPI.CreatePaymail(ctx, &commands.CreatePaymail{
		Metadata: queryparams.Metadata{"key": "value"},
		Key:      keys.XPub(),
		Address:  examples.Paymail,
	})
	if err != nil {
		log.Fatalf("Failed to create paymail: %v", err)
	}
	exampleutil.PrettyPrint("Created paymail", paymail)
}
