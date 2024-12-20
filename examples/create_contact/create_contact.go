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
	ctx := context.Background()
	adminAPI, err := wallet.NewAdminAPIWithXPriv(exampleutil.NewDefaultConfig(), examples.XPriv)
	if err != nil {
		log.Fatal(err)
	}

	contact, err := adminAPI.CreateContact(ctx, &commands.CreateContact{
		CreatorPaymail: "john.doe@example.com",
		FullName:       "Jane Doe",
	})
	if err != nil {
		log.Fatal(err)
	}

	exampleutil.Print("Create Paymail - api/v1/admin/contacts/paymail", contact)
}
