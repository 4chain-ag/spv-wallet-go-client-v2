package auth_test

import (
	"net/http"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/require"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/auth"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/spvwallettest"
)

const (
	xAuthKey          = "X-Auth-Key"
	xAuthXPubKey      = "X-Auth-Xpub"
	xAuthHashKey      = "X-Auth-Hash"
	xAuthNonceKey     = "X-Auth-Nonce"
	xAuthTimeKey      = "X-Auth-Time"
	xAuthSignatureKey = "X-Auth-Signature"
)

// func TestAccessKeyAuthenitcator_NewXprivAuthenticatorFromString(t *testing.T) {
// 	// Given:
// 	xpriv := spvwallettest.ExtendedKey(t)
// 	authenticator, err := auth.NewXprivAuthenticator(xpriv)

// 	// When:
// 	authenticatorFromStr, err := auth.NewXprivAuthenticatorFromString(spvwallettest.UserXPriv)

// 	// Then:
// 	require.NotNil(t, authenticator)
// 	require.NoError(t, err)
// 	require.Equal(t, authenticator, authenticatorFromStr)
// }

// func TestAccessKeyAuthenitcator_NewAccessKeyAuthenticatorFromString(t *testing.T) {
// 	// Given:
// 	privKeyStr := spvwallettest.PrivateKeyHexString(t)

// 	// When:
// 	authenticator, err := auth.NewAccessKeyAuthenticatorFromString(privKeyStr)

// 	// Then:
// 	require.NotNil(t, authenticator)
// 	require.NoError(t, err)
// }

// func TestAccessKeyAuthenitcator_NewXpubOnlyAuthenticatorFromString(t *testing.T) {
// 	// Given:
// 	xpubStr := spvwallettest.ExtendedKeyString(t)

// 	// When:
// 	authenticator, err := auth.NewXpubOnlyAuthenticatorFromString(xpubStr)

// 	// Then:
// 	require.NotNil(t, authenticator)
// 	require.NoError(t, err)
// }

func TestAccessKeyAuthenitcator_NewWithNilAccessKey(t *testing.T) {
	// when:
	authenticator, err := auth.NewAccessKeyAuthenticator(nil)

	// then:
	require.Nil(t, authenticator)
	require.ErrorIs(t, err, auth.ErrEcPrivateKey)
}

func TestAccessKeyAuthenticator_Authenticate(t *testing.T) {
	// given:
	key := spvwallettest.PrivateKey(t)
	authenticator, err := auth.NewAccessKeyAuthenticator(key)
	require.NotNil(t, authenticator)
	require.NoError(t, err)

	req := resty.New().R()

	// when:
	err = authenticator.Authenticate(req)

	// then:
	require.NoError(t, err)
	requireXAuthHeaderToBeSet(t, req.Header)
	requireSignatureHeadersToBeSet(t, req.Header)
}

func TestXprivAuthenitcator_NewWithNilXpriv(t *testing.T) {
	// when:
	authenticator, err := auth.NewXprivAuthenticator("")

	// then:
	require.Nil(t, authenticator)
	require.ErrorIs(t, err, auth.ErrBip32ExtendedKey)
}

func TestXprivAuthenitcator_Authenticate(t *testing.T) {
	// given:
	authenticator, err := auth.NewXprivAuthenticator(spvwallettest.UserXPriv)
	require.NotNil(t, authenticator)
	require.NoError(t, err)

	req := resty.New().R()

	// when:
	err = authenticator.Authenticate(req)

	// then:
	require.NoError(t, err)
	requireXpubHeaderToBeSet(t, req.Header)
	requireSignatureHeadersToBeSet(t, req.Header)
}

func TestXpubOnlyAuthenticator_NewWithNilXpub(t *testing.T) {
	// when:
	authenticator, err := auth.NewXpubOnlyAuthenticator(nil)

	// then:
	require.Nil(t, authenticator)
	require.ErrorIs(t, err, auth.ErrBip32ExtendedKey)
}

func TestXpubOnlyAuthenticator_Authenticate(t *testing.T) {
	// given:
	key := spvwallettest.ExtendedKey(t)

	authenticator, err := auth.NewXpubOnlyAuthenticator(key)
	require.NotNil(t, authenticator)
	require.NoError(t, err)

	req := resty.New().R()

	// when:
	err = authenticator.Authenticate(req)

	// then:
	require.NoError(t, err)
	requireXpubHeaderToBeSet(t, req.Header)
}

func requireXAuthHeaderToBeSet(t *testing.T, h http.Header) {
	require.Equal(t, []string{spvwallettest.UserPubAccessKey}, h[xAuthKey])
}

func requireXpubHeaderToBeSet(t *testing.T, h http.Header) {
	require.Equal(t, []string{spvwallettest.UserXPub}, h[xAuthXPubKey])
}

func requireSignatureHeadersToBeSet(t *testing.T, h http.Header) {
	expected := []string{
		xAuthHashKey,
		xAuthNonceKey,
		xAuthTimeKey,
		xAuthSignatureKey,
	}

	actual := make([]string, 0, len(expected))
	for k := range h {
		actual = append(actual, k)
	}
	require.Subset(t, actual, expected)
}
