package main

import (
	"context"
	"log"

	wallet "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples/exampleutil"
)

func main() {
	usersAPI, err := wallet.NewUserAPIWithXPriv(exampleutil.NewDefaultConfig(), examples.XPriv)
	if err != nil {
		log.Fatal(err)
	}

	paymail := "john.doe@example.com"
	err = usersAPI.RemoveContact(context.Background(), paymail)
	if err != nil {
		log.Fatal(err)
	}

	exampleutil.Print("\n[HTTP DELETE] Remove contact - api/v1/contacts/%s", paymail)
}
