package main

import (
	"fmt"
	"log"

	"github.com/bitcoin-sv/spv-wallet-go-client/docs/keysgen"
)

func main() {
	keys, err := keysgen.GenerateKeys()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("XPriv: ", keys.Xpriv)
	fmt.Println("XPub: ", keys.Xpub)
	fmt.Println("Mnemonic: ", keys.Mnemonic)
}
