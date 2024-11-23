package walletkeys

import (
	"fmt"

	bip32 "github.com/bitcoin-sv/go-sdk/compat/bip32"
	bip39 "github.com/bitcoin-sv/go-sdk/compat/bip39"
	chaincfg "github.com/bitcoin-sv/go-sdk/transaction/chaincfg"
)

const DefaultEntropy = 128 // The bit size must be a multiple of 32 and fall within the inclusive range of {128, 256}.

// Keys represents a set of hierarchical deterministic (HD) keys,
// including the extended private key (XPriv) and extended public key (XPub).
type Keys struct {
	xPriv string
	xPub  string
}

// XPriv returns the HD extended private key as a string.
func (k *Keys) XPriv() string { return k.xPriv }

// XPub returns the HD extended public key as a string.
func (k *Keys) XPub() string { return k.xPub }

// KeysWithMnemonic extends the Keys struct by including the mnemonic phrase
// used to generate the associated xPriv and XPub HD keys as strings.
type KeysWithMnemonic struct {
	keys     Keys
	mnemonic string
}

// Mnemonic returns the mnemonic phrase used to generate the keys.
func (k *KeysWithMnemonic) Mnemonic() string { return k.mnemonic }

// Keys returns the associated Keys struct.
func (k *KeysWithMnemonic) Keys() Keys { return k.keys }

// XPrivFromString generates an extended private key (xPriv) from a string.
// It returns the extended private key and an error if the conversion fails.
func XPrivFromString(s string) (*bip32.ExtendedKey, error) {
	xPriv, err := bip32.NewKeyFromString(s)
	if err != nil {
		return nil, fmt.Errorf("failed to generate HD key from string: %w", err)
	}

	return xPriv, nil
}

// XPrivFromMnemonic generates an extended private key (xPriv) from a mnemonic phrase.
// It returns the extended private key and an error if seed generation or HD key creation fails.
func XPrivFromMnemonic(mnemonic string) (*bip32.ExtendedKey, error) {
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return nil, fmt.Errorf("failed to generate seed from mnemonic: %w", err)
	}

	xPriv, err := bip32.NewMaster(seed, &chaincfg.MainNet)
	if err != nil {
		return nil, fmt.Errorf("failed to create master node HD key: %w", err)
	}

	return xPriv, nil
}

// RandomXPriv generates a random extended private key (xPriv).
// The seed size is specified as 32 bytes (256 bits), as defined by the bip32.RecommendedSeedLen constant.
// It returns a pointer to the extended private key and an error if seed generation or the creation of the master node HD key fails.
func RandomXPriv() (*bip32.ExtendedKey, error) {
	seed, err := bip32.GenerateSeed(bip32.RecommendedSeedLen)
	if err != nil {
		return nil, fmt.Errorf("failed to generate seed: %w", err)
	}

	xPriv, err := bip32.NewMaster(seed, &chaincfg.MainNet)
	if err != nil {
		return nil, fmt.Errorf("failed to generate master node HD key: %w", err)
	}

	return xPriv, nil
}

// RandomMnemonic generates a mnemonic phrase consisting of words derived from default entropy.
// It returns the mnemonic as a string and an error if entropy generation or mnemonic creation fails.
func RandomMnemonic() (string, error) {
	entropy, err := bip39.NewEntropy(DefaultEntropy)
	if err != nil {
		return "", fmt.Errorf("failed to generate entropy: %w", err)
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", fmt.Errorf("failed to generate mnemonic: %w", err)
	}

	return mnemonic, nil
}

// RandomKeys generates random HD keys (xPriv and xPub).
// It returns a Keys struct containing the extended private and public keys and an error if any generation fails.
func RandomKeys() (*Keys, error) {
	xPriv, err := RandomXPriv()
	if err != nil {
		return nil, fmt.Errorf("failed to generate random xPriv: %w", err)
	}

	xPub, err := bip32.GetExtendedPublicKey(xPriv)
	if err != nil {
		return nil, fmt.Errorf("failed to get extended public key: %w", err)
	}

	return &Keys{xPriv: xPriv.String(), xPub: xPub}, nil
}

// RandomKeysWithMnemonic generates random HD keys (xPriv and xPub) along with a mnemonic phrase.
// It returns a KeysWithMnemonic struct containing the keys and the associated mnemonic, and an error if any generation fails.
func RandomKeysWithMnemonic() (*KeysWithMnemonic, error) {
	mnemonic, err := RandomMnemonic()
	if err != nil {
		return nil, fmt.Errorf("failed to generate random mnemonic: %w", err)
	}

	xPriv, err := bip32.GenerateHDKeyFromMnemonic(mnemonic, "", &chaincfg.MainNet)
	if err != nil {
		return nil, fmt.Errorf("failed to generate HD key from mnemonic: %w", err)
	}

	xPub, err := bip32.GetExtendedPublicKey(xPriv)
	if err != nil {
		return nil, fmt.Errorf("failed to get extended public key: %w", err)
	}

	keys := Keys{xPriv: xPriv.String(), xPub: xPub}
	return &KeysWithMnemonic{mnemonic: mnemonic, keys: keys}, nil
}
