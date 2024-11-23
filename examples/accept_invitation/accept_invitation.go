package main

import (
	"context"
	"fmt"
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

	paymail := "john.doe@example.com"
	err = spv.AcceptInvitation(context.Background(), paymail)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(fmt.Sprintf("\n[HTTP POST] Accept contact invitation - api/v1/invitations/%s/contacts", paymail))
}
