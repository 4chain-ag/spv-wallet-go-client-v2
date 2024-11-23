package main

import (
	"fmt"
	"log"

	"github.com/bitcoin-sv/spv-wallet-go-client/walletkeys"
)

func main() {
	keys, err := walletkeys.RandomKeysWithMnemonic()
	if err != nil {
		log.Fatal(err)
	}

	inner := keys.Keys()
	fmt.Println("XPriv: ", inner.XPriv())
	fmt.Println("XPub: ", inner.XPub())
	fmt.Println("Mnemonic: ", keys.Mnemonic())
}
