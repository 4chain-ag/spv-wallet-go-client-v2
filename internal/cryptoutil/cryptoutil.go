package cryptoutil

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"

	bip32 "github.com/bitcoin-sv/go-sdk/compat/bip32"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
)

const (
	// XpubKeyLength is the length of an xPub string key
	XpubKeyLength = 111

	// ChainInternal internal chain num
	ChainInternal = uint32(1)

	// ChainExternal external chain num
	ChainExternal = uint32(0)

	// MaxInt32 max integer for int32
	MaxInt32 = int64(1<<(32-1) - 1)
)

// Hash returns the sha256 hash of the data string
func Hash(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// RandomHex returns a random hex string and error
func RandomHex(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// DeriveChildKeyFromHex derive the child extended key from the hex string
func DeriveChildKeyFromHex(hdKey *bip32.ExtendedKey, hexHash string) (*bip32.ExtendedKey, error) {
	var childKey *bip32.ExtendedKey
	childKey = hdKey

	childNums, err := utils.GetChildNumsFromHex(hexHash)
	if err != nil {
		return nil, err
	}

	for _, num := range childNums {
		if childKey, err = childKey.Child(num); err != nil {
			return nil, err
		}
	}
	return childKey, nil
}
