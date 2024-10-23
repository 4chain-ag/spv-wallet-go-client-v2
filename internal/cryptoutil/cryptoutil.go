package cryptoutil

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"strconv"

	bip32 "github.com/bitcoin-sv/go-sdk/compat/bip32"
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

	childNums, err := GetChildNumsFromHex(hexHash)
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

// GetChildNumsFromHex get an array of uint32 numbers from the hex string
func GetChildNumsFromHex(hexHash string) ([]uint32, error) {
	strLen := len(hexHash)
	size := 8
	splitLength := int(math.Ceil(float64(strLen) / float64(size)))
	childNums := make([]uint32, 0)
	for i := 0; i < splitLength; i++ {
		start := i * size
		stop := start + size
		if stop > strLen {
			stop = strLen
		}
		num, err := strconv.ParseInt(hexHash[start:stop], 16, 64)
		if err != nil {
			return nil, fmt.Errorf("parse int op failure: %w", err)
		}

		result := num % MaxInt32
		resultU32, err := Int64ToUint32(result)
		if err != nil {
			return nil, fmt.Errorf("int64 to uint32 convert op failure: %w", err)
		}
		childNums = append(childNums, resultU32)
	}

	return childNums, nil
}

// Int64ToUint32 converts an int64 value to uint32 with range checks.
// It returns the converted uint32 value and a nil error if the conversion is successful.
// Otherwise, it returns the zero value and a non-nil error that describes the reason for the conversion failure.
func Int64ToUint32(value int64) (uint32, error) {
	if value < 0 {
		return 0, ErrNegativeValueNotAllowed
	}
	if value > math.MaxUint32 {
		return 0, ErrMaxUint32LimitExceeded
	}
	return uint32(value), nil
}

var (
	// ErrMaxUint32LimitExceeded occurs when attempting to convert an int64 value that exceeds the maximum uint32 limit.
	ErrMaxUint32LimitExceeded = errors.New("max uint32 value exceeded")

	// ErrNegativeValueNotAllowed occurs when attempting to convert a negative int64 value to uint32.
	ErrNegativeValueNotAllowed = errors.New("negative value is not allowed")
)
