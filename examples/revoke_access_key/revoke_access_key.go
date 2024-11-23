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

	accessKeyID := "9f7efc4af8f2c9f745ca8dfa737394d810dd8828c072c7c05e07c7aae67ff790"
	err = spv.RevokeAccessKey(context.Background(), accessKeyID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(fmt.Sprintf("\n[HTTP DELETE] Revoke access key - api/v1/users/current/keys/%s", accessKeyID))
}
