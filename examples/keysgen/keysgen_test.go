package keysgen_test

import (
	"fmt"
	"log"

	"github.com/bitcoin-sv/spv-wallet-go-client/examples/keysgen"
)

func ExampleGenerateKeys() {
	keys, err := keysgen.GenerateKeys()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("xPriv: ", keys.Xpriv)
	fmt.Println("Xpub: ", keys.Xpub)
}

func ExampleGenerateKeysFromString() {
	key := "xprv9s21ZrQH143K3Lh4wdicqvYNMcdh49rMLqDvQoyys8L6f5tfE2WkQN7ZVE2awBrfVWNSJ8pPd4QLLr94Nur85Dvj8kD8RoZghBuNTpvL8si"
	keys, err := keysgen.GenerateKeysFromString(key)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("xpub:", keys.Xpub)
	fmt.Println("xpriv:", keys.Xpriv)

	// Output:
	// xpub: xpub661MyMwAqRbcFpmY3fFdD4V6ueUBTcaCi49XDCPbRTs5XtDomZpzxAS3LUb2hMfUVphDsSPxfjietmsBRFkLDY9Xa3P4jbgNDMnDK3UqJe2
	// xpriv: xprv9s21ZrQH143K3Lh4wdicqvYNMcdh49rMLqDvQoyys8L6f5tfE2WkQN7ZVE2awBrfVWNSJ8pPd4QLLr94Nur85Dvj8kD8RoZghBuNTpvL8si
}

func ExampleGenerateKeysFromMnemonic() {
	mnemonic := "absorb corn ostrich order sing boost just harvest enable make detail future desert bus adult"
	keys, err := keysgen.GenerateKeysFromMnemonic(mnemonic)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("mnemonic:", keys.Mnemonic)
	fmt.Println("xpub:", keys.Xpub)
	fmt.Println("xpriv:", keys.Xpriv)

	// Output:
	// mnemonic: absorb corn ostrich order sing boost just harvest enable make detail future desert bus adult
	// xpub: xpub661MyMwAqRbcFpmY3fFdD4V6ueUBTcaCi49XDCPbRTs5XtDomZpzxAS3LUb2hMfUVphDsSPxfjietmsBRFkLDY9Xa3P4jbgNDMnDK3UqJe2
	// xpriv: xprv9s21ZrQH143K3Lh4wdicqvYNMcdh49rMLqDvQoyys8L6f5tfE2WkQN7ZVE2awBrfVWNSJ8pPd4QLLr94Nur85Dvj8kD8RoZghBuNTpvL8si
}
