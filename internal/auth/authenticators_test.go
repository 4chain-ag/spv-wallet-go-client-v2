package auth_test

import (
	"encoding/hex"
	"net/http"
	"testing"

	bip32 "github.com/bitcoin-sv/go-sdk/compat/bip32"
	compat "github.com/bitcoin-sv/go-sdk/compat/bip32"
	ec "github.com/bitcoin-sv/go-sdk/primitives/ec"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/auth"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/require"
)

func TestXpubAuthenticator_Authenticate(t *testing.T) {
	key := extendedKey(t)
	authenticator, err := auth.NewXpubOnlyAuthenticator(key)
	require.NotNil(t, authenticator)
	require.NoError(t, err)

	req := resty.New().R()
	err = authenticator.Authenticate(req)
	require.NoError(t, err)
	xPubHeadersTestHelper(t, req.Header, key)
}

func TestAccessKeyAuthenitcator_Authenticate(t *testing.T) {
	key := privateKey(t)
	authenticator, err := auth.NewAccessKeyAuthenticator(key)
	require.NotNil(t, authenticator)
	require.NoError(t, err)

	req := resty.New().R()
	err = authenticator.Authenticate(req)
	require.NoError(t, err)
	accessKeyHeadersTestHelper(t, req.Header, hex.EncodeToString(key.PubKey().SerializeCompressed()))
}

func TestXprivAuthenitcator_Authenticate(t *testing.T) {
	key := extendedKey(t)
	authenticator, err := auth.NewXprivAuthenticator(key)
	require.NotNil(t, authenticator)
	require.NoError(t, err)

	req := resty.New().R()
	err = authenticator.Authenticate(req)
	require.NoError(t, err)
	xPrivHeadersTestHelper(t, req.Header, key)
}

func extendedKey(t *testing.T) *bip32.ExtendedKey {
	t.Helper()
	key, err := bip32.GenerateHDKey(compat.RecommendedSeedLength)
	if err != nil {
		t.Fatalf("test helper - failed to generate HD key from string: %sw", err)
	}
	return key
}

func privateKey(t *testing.T) *ec.PrivateKey {
	t.Helper()
	key, err := ec.NewPrivateKey()
	if err != nil {
		t.Fatalf("test helper - failed to create private key: %s", err)
	}
	return key
}

func xPubHeadersTestHelper(t *testing.T, h http.Header, key *compat.ExtendedKey) {
	t.Helper()
	xPub, err := key.Neuter()
	if err != nil {
		t.Fatalf("test helper - failed to get extended public key: %s", err)
	}

	actualXpub := h["X-Auth-Xpub"]
	expectedXpub := []string{xPub.String()}
	require.Equal(t, expectedXpub, actualXpub)
}

func xPrivHeadersTestHelper(t *testing.T, h http.Header, key *compat.ExtendedKey) {
	t.Helper()
	xPub, err := key.Neuter()
	if err != nil {
		t.Fatalf("test helper - failed to get extended public key: %s", err)
	}

	expectedXpub := []string{xPub.String()}
	expectedHeaders := []string{
		"X-Auth-Xpub",
		"X-Auth-Hash",
		"X-Auth-Nonce",
		"X-Auth-Time",
		"X-Auth-Signature",
	}
	actualXpub := h["X-Auth-Xpub"]
	actualHeaders := make([]string, 0, len(expectedHeaders))
	for k := range h {
		actualHeaders = append(actualHeaders, k)
	}

	require.ElementsMatch(t, expectedHeaders, actualHeaders)
	require.Equal(t, expectedXpub, actualXpub)
}

func accessKeyHeadersTestHelper(t *testing.T, h http.Header, key string) {
	t.Helper()
	expectedAuthKey := []string{key}
	expectedHeaders := []string{
		"X-Auth-Key",
		"X-Auth-Hash",
		"X-Auth-Nonce",
		"X-Auth-Time",
		"X-Auth-Signature",
	}
	actualAuthKey := h["X-Auth-Key"]
	actualHeaders := make([]string, 0, len(expectedHeaders))
	for k := range h {
		actualHeaders = append(actualHeaders, k)
	}

	require.ElementsMatch(t, expectedHeaders, actualHeaders)
	require.Equal(t, expectedAuthKey, actualAuthKey)
}
