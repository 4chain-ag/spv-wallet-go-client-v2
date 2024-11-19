package keysgen

import (
	"fmt"

	bip32 "github.com/bitcoin-sv/go-sdk/compat/bip32"
	bip39 "github.com/bitcoin-sv/go-sdk/compat/bip39"
	chaincfg "github.com/bitcoin-sv/go-sdk/transaction/chaincfg"
)

// Keys represents a set of hierarchical deterministic (HD) keys,
// including the extended private key (Xpriv), extended public key (Xpub),
// and the mnemonic phrase used to generate these keys.
type Keys struct {
	Xpriv    string // The HD extended private key as a string.
	Xpub     string // The HD extended public key as a string.
	Mnemonic string // The mnemonic phrase used to derive the keys, if available.
}

// GenerateKeysFromString creates a new set of HD keys from an extended private key string.
func GenerateKeysFromString(s string) (*Keys, error) {
	xPriv, err := bip32.NewKeyFromString(s)
	if err != nil {
		return nil, fmt.Errorf("failed to create HD private key: %w", err)
	}

	xPub, err := xPriv.Neuter()
	if err != nil {
		return nil, fmt.Errorf("failed to create HD public key: %w", err)
	}

	return &Keys{Xpriv: xPriv.String(), Xpub: xPub.String()}, nil
}

// GenerateKeysFromMnemonic creates a new set of HD keys from a mnemonic phrase.
func GenerateKeysFromMnemonic(mnemonic string) (*Keys, error) {
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return nil, fmt.Errorf("failed to create seed: %w", err)
	}

	xPriv, err := bip32.NewMaster(seed, &chaincfg.MainNet)
	if err != nil {
		return nil, fmt.Errorf("failed to create HD private key: %w", err)
	}

	xPub, err := xPriv.Neuter()
	if err != nil {
		return nil, fmt.Errorf("failed to create HD public key: %w", err)
	}

	return &Keys{Xpriv: xPriv.String(), Xpub: xPub.String(), Mnemonic: mnemonic}, nil
}

// GenerateKeys generates a new set of HD keys using a randomly generated mnemonic phrase.
func GenerateKeys() (*Keys, error) {
	entropy, err := bip39.NewEntropy(160)
	if err != nil {
		return nil, fmt.Errorf("failed to generate entropy: %w", err)
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return nil, fmt.Errorf("failed to generate mnemonic: %w", err)
	}

	hdKey, err := bip32.GenerateHDKeyFromMnemonic(mnemonic, "", &chaincfg.MainNet)
	if err != nil {
		return nil, fmt.Errorf("failed to generate HD key from mnemonic: %w", err)
	}

	xPriv := hdKey.String()
	xPub, err := bip32.GetExtendedPublicKey(hdKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get extended public key from HD key: %w", err)
	}

	return &Keys{Xpriv: xPriv, Xpub: xPub, Mnemonic: mnemonic}, nil
}
