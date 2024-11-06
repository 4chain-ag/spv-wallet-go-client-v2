package main

import (
	"fmt"
	"log"

	bip32 "github.com/bitcoin-sv/go-sdk/compat/bip32"
	bip39 "github.com/bitcoin-sv/go-sdk/compat/bip39"
	chaincfg "github.com/bitcoin-sv/go-sdk/transaction/chaincfg"
)

func main() {
	keys, err := GenerateKeys()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("xpriv: ", keys.xpriv)
	fmt.Println("xpub: ", keys.xpub)
	fmt.Println("mnemonic: ", keys.mnemonic)
}

type Keys struct {
	xpriv    string
	xpub     string
	mnemonic string
}

func GenerateKeys() (*Keys, error) {
	entropy, err := bip39.NewEntropy(160)
	if err != nil {
		return nil, fmt.Errorf("failed to generate entropy: %w", err)
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return nil, fmt.Errorf("failed to generate generate mnemonic: %w", err)
	}

	hdKey, err := bip32.GenerateHDKeyFromMnemonic(mnemonic, "", &chaincfg.MainNet)
	if err != nil {
		return nil, fmt.Errorf("failed to generate hd key from mnemonic: %w", err)
	}

	hdXpriv := hdKey.String()
	hdXpub, err := bip32.GetExtendedPublicKey(hdKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get extended public key from hd key: %w", err)
	}

	keys := &Keys{
		xpriv:    hdXpriv,
		xpub:     hdXpub,
		mnemonic: mnemonic,
	}
	return keys, nil
}
