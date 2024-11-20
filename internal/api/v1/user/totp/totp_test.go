package totp_test

import (
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	bip32 "github.com/bitcoin-sv/go-sdk/compat/bip32"
	client "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/totp"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/xpriv"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/clienttest"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/stretchr/testify/require"
)

func TestClient_GenerateTotpForContact(t *testing.T) {
	cfg := client.Config{
		Addr:    clienttest.TestAPIAddr,
		Timeout: 5 * time.Second,
	}
	t.Run("success", func(t *testing.T) {
		// given
		sut, err := client.NewWithXPriv(cfg, clienttest.UserXPriv)
		require.NoError(t, err)
		require.NotNil(t, sut.ClientXPriv)

		contact := models.Contact{PubKey: clienttest.PubKey}
		wc := totp.New(sut)
		// when
		pass, err := wc.GenerateTotpForContact(&contact, 30, 2)

		// then
		require.NoError(t, err)
		require.Len(t, pass, 2)
	})

	t.Run("WalletClient without xPriv - returns error", func(t *testing.T) {
		// given
		sut, err := client.NewWithXPub(cfg, clienttest.UserXPub)
		require.NoError(t, err)
		require.NotNil(t, sut.XPub)
		wc := totp.New(sut)
		// when
		_, err = wc.GenerateTotpForContact(nil, 30, 2)

		// then
		require.ErrorIs(t, err, totp.ErrMissingXpriv)
	})

	t.Run("contact has invalid PubKey - returns error", func(t *testing.T) {
		// given
		sut, err := client.NewWithXPriv(cfg, clienttest.UserXPriv)
		require.NoError(t, err)
		require.NotNil(t, sut.ClientXPriv)

		contact := models.Contact{PubKey: "invalid-pk-format"}
		wc := totp.New(sut)
		// when
		_, err = wc.GenerateTotpForContact(&contact, 30, 2)

		// then
		require.ErrorIs(t, err, totp.ErrContactPubKeyInvalid)

	})
}

func TestClient_ValidateTotpForContact(t *testing.T) {
	cfg := client.Config{
		Addr:    clienttest.TestAPIAddr,
		Timeout: 5 * time.Second,
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This handler could be adjusted depending on the expected API endpoints
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("123456")) // Simulate a TOTP response for any requests
	}))
	defer server.Close()
	t.Run("success", func(t *testing.T) {
		aliceKeys, err := xpriv.Generate()
		require.NoError(t, err)
		bobKeys, err := xpriv.Generate()
		require.NoError(t, err)

		// Set up the WalletClient for Alice and Bob
		clientAlice, err := client.NewWithXPriv(cfg, aliceKeys.XPriv())
		require.NoError(t, err)
		require.NotNil(t, clientAlice.ClientXPriv)
		clientBob, err := client.NewWithXPriv(cfg, bobKeys.XPriv())
		require.NoError(t, err)
		require.NotNil(t, clientBob.ClientXPriv)

		aliceContact := &models.Contact{
			PubKey:  makeMockPKI(aliceKeys.XPub().String()),
			Paymail: "bob@example.com",
		}

		bobContact := &models.Contact{
			PubKey:  makeMockPKI(bobKeys.XPub().String()),
			Paymail: "bob@example.com",
		}

		wcAlice := totp.New(clientAlice)
		wcBob := totp.New(clientBob)
		// Generate and validate TOTP
		passcode, err := wcAlice.GenerateTotpForContact(bobContact, 3600, 6)
		require.NoError(t, err)
		result, err := wcBob.ValidateTotpForContact(aliceContact, passcode, bobContact.Paymail, 3600, 6)
		require.NoError(t, err)
		require.True(t, result)
	})

	t.Run("contact has invalid PubKey - returns error", func(t *testing.T) {
		sut, err := client.NewWithXPriv(cfg, clienttest.UserXPriv)
		require.NoError(t, err)

		invalidContact := &models.Contact{
			PubKey:  "invalid_pub_key_format",
			Paymail: "invalid@example.com",
		}
		wc := totp.New(sut)
		_, err = wc.ValidateTotpForContact(invalidContact, "123456", "someone@example.com", 3600, 6)
		require.Error(t, err)
		require.Contains(t, err.Error(), "contact's PubKey is invalid")
	})
}

func makeMockPKI(xpub string) string {
	xPub, _ := bip32.NewKeyFromString(xpub)
	var err error
	for i := 0; i < 3; i++ { //magicNumberOfInheritance is 3 -> 2+1; 2: because of the way spv-wallet stores xpubs in db; 1: to make a PKI
		xPub, err = xPub.Child(0)
		if err != nil {
			panic(err)
		}
	}

	pubKey, err := xPub.ECPubKey()
	if err != nil {
		panic(err)
	}

	return hex.EncodeToString(pubKey.SerializeCompressed())
}
