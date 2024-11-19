package totp

import (
	"encoding/base32"
	"encoding/hex"
	"time"

	bip32 "github.com/bitcoin-sv/go-sdk/compat/bip32"
	ec "github.com/bitcoin-sv/go-sdk/primitives/ec"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"

	client "github.com/bitcoin-sv/spv-wallet-go-client"
	utils "github.com/bitcoin-sv/spv-wallet-go-client/internal/cryptoutil"
)

const (
	// TotpDefaultPeriod - Default number of seconds a TOTP is valid for.
	TotpDefaultPeriod uint = 30
	// TotpDefaultDigits - Default TOTP length
	TotpDefaultDigits uint = 2
)

// WalletClient handles TOTP operations for SPV Wallet.
type WalletClient struct {
	client *client.Client
}

// New creates a new TOTP WalletClient.
func New(c *client.Client) *WalletClient {
	return &WalletClient{client: c}
}

// GenerateTotpForContact generates a time-based one-time password (TOTP) for a contact.
func (b *WalletClient) GenerateTotpForContact(contact *models.Contact, period, digits uint) (string, error) {
	sharedSecret, err := makeSharedSecret(b.client, contact)
	if err != nil {
		return "", err
	}

	opts := getTotpOpts(period, digits)
	return totp.GenerateCodeCustom(directedSecret(sharedSecret, contact.Paymail), time.Now(), *opts)
}

// ValidateTotpForContact validates a TOTP for a contact.
func (b *WalletClient) ValidateTotpForContact(contact *models.Contact, passcode, requesterPaymail string, period, digits uint) (bool, error) {
	sharedSecret, err := makeSharedSecret(b.client, contact)
	if err != nil {
		return false, err
	}

	opts := getTotpOpts(period, digits)
	return totp.ValidateCustom(passcode, directedSecret(sharedSecret, requesterPaymail), time.Now(), *opts)
}

func makeSharedSecret(client *client.Client, contact *models.Contact) ([]byte, error) {
	privKey, pubKey, err := getSharedSecretFactors(client, contact)
	if err != nil {
		return nil, err
	}

	x, _ := ec.S256().ScalarMult(pubKey.X, pubKey.Y, privKey.D.Bytes())
	return x.Bytes(), nil
}

func getSharedSecretFactors(client *client.Client, contact *models.Contact) (*ec.PrivateKey, *ec.PublicKey, error) {
	// Retrieve xPriv from client or configuration.
	xPriv := client.XPriv()
	if xPriv == nil {
		return nil, nil, ErrMissingXpriv
	}

	// Derive private key from xPriv for PKI operations.
	xpriv, err := deriveXprivForPki(xPriv)
	if err != nil {
		return nil, nil, err
	}

	privKey, err := xpriv.ECPrivKey()
	if err != nil {
		return nil, nil, err
	}

	// Convert contact's public key.
	pubKey, err := convertPubKey(contact.PubKey)
	if err != nil {
		return nil, nil, ErrContactPubKeyInvalid
	}

	return privKey, pubKey, nil
}

func deriveXprivForPki(xpriv *bip32.ExtendedKey) (*bip32.ExtendedKey, error) {
	pkiXpriv, err := bip32.GetHDKeyByPath(xpriv, utils.ChainExternal, 0)
	if err != nil {
		return nil, err
	}

	return pkiXpriv.Child(0)
}

func convertPubKey(pubKey string) (*ec.PublicKey, error) {
	decoded, err := hex.DecodeString(pubKey)
	if err != nil {
		return nil, err
	}

	return ec.ParsePubKey(decoded)
}

func getTotpOpts(period, digits uint) *totp.ValidateOpts {
	if period == 0 {
		period = TotpDefaultPeriod
	}

	if digits == 0 {
		digits = TotpDefaultDigits
	}

	return &totp.ValidateOpts{
		Period: period,
		Digits: otp.Digits(digits),
	}
}

// directedSecret appends a paymail to the shared secret and encodes it as base32.
func directedSecret(sharedSecret []byte, paymail string) string {
	return base32.StdEncoding.EncodeToString(append(sharedSecret, []byte(paymail)...))
}
