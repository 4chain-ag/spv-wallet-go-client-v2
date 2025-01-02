package main

import (
	"fmt"
	"log"

	wallet "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples"
	"github.com/bitcoin-sv/spv-wallet-go-client/examples/exampleutil"
	"github.com/bitcoin-sv/spv-wallet/models"
)

func main() {
	const aliceXPriv = "xprv9s21ZrQH143K4JFXqGhBzdrthyNFNuHPaMUwvuo8xvpHwWXprNK7T4JPj1w53S1gojQncyj8JhSh8qouYPZpbocsq934cH5G1t1DRBfgbod"

	// pubKey - PKI can be obtained from the contact's paymail capability
	const bobPKI = "03a48e13dc598dce5fda9b14ea13f32d5dbc4e8d8a34447dda84f9f4c457d57fe7"
	const digits = 4
	const period = 1200

	usersAPI, err := wallet.NewUserAPIWithXPriv(exampleutil.NewDefaultConfig(), examples.UserXPriv)
	if err != nil {
		log.Fatalf("Failed to initialize user API with XPriv: %v", err)
	}

	mockContact := &models.Contact{
		PubKey:  bobPKI,
		Paymail: "test@paymail.com",
	}
	code, err := usersAPI.GenerateTotpForContact(mockContact, period, digits)
	if err != nil {
		log.Fatalf("Failed to generate totp for contact: %v", err)
	}

	fmt.Println("TOTP code from Alice to Bob: ", code)

	err = usersAPI.ValidateTotpForContact(mockContact, code, mockContact.Paymail, period, digits)
	if err != nil {
		log.Fatalf("Failed to validate totp for contact: %v", err)
	}

	fmt.Println("TOTP code from Alice to Bob is valid")
}
