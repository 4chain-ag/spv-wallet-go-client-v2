package auth

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"

	bip32 "github.com/bitcoin-sv/go-sdk/compat/bip32"
	ec "github.com/bitcoin-sv/go-sdk/primitives/ec"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/go-resty/resty/v2"
)

// XpubAuthenticator provides authentication using an extended public key (xpub).
type XpubAuthenticator struct {
	hdKey *bip32.ExtendedKey
}

// Authenticate sets an xpub-based authentication header on the given HTTP request.
func (x *XpubAuthenticator) Authenticate(r *resty.Request) error {
	xPub, err := bip32.GetExtendedPublicKey(x.hdKey)
	if err != nil {
		return fmt.Errorf("failed to get extended public key: %w", err)
	}
	r.SetHeader(models.AuthHeader, xPub)
	return nil
}

// XprivAuthenticator provides authentication using an extended private key (xpriv).
type XprivAuthenticator struct {
	xpubAuth *XpubAuthenticator
	xpriv    *bip32.ExtendedKey
}

// Authenticate sets an xpub header and adds a signature generated from the xpriv.
func (x *XprivAuthenticator) Authenticate(r *resty.Request) error {
	err := x.xpubAuth.Authenticate(r)
	if err != nil {
		return fmt.Errorf("failed to set xpub header: %w", err)
	}
	body := bodyString(r)
	header := make(http.Header)
	err = setSignature(&header, x.xpriv, body)
	if err != nil {
		return fmt.Errorf("failed to sign request with xpriv: %w", err)
	}
	r.SetHeaderMultiValues(header)
	return nil
}

// AccessKeyAuthenticator provides authentication using EC private/public keys.
type AccessKeyAuthenticator struct {
	priv *ec.PrivateKey
	pub  *ec.PublicKey
}

// Authenticate sets an access key header and a signature on the request.
func (a *AccessKeyAuthenticator) Authenticate(r *resty.Request) error {
	r.Header.Set(models.AuthAccessKey, a.pubKeyHex())
	body := bodyString(r)
	sign, err := createSignatureAccessKey(a.privKeyHex(), body)
	if err != nil {
		return fmt.Errorf("failed to sign request with access key: %w", err)
	}
	setSignatureHeaders(&r.Header, sign)
	return nil
}

func (a *AccessKeyAuthenticator) privKeyHex() string {
	return hex.EncodeToString(a.priv.Serialize())
}

func (a *AccessKeyAuthenticator) pubKeyHex() string {
	return hex.EncodeToString(a.pub.SerializeCompressed())
}

func bodyString(r *resty.Request) string {
	switch r.Method {
	case http.MethodGet:
		return ""
	}
	return ""
}

// NewXprivAuthenticator creates an XprivAuthenticator with the provided xpriv key.
func NewXprivAuthenticator(xpriv *bip32.ExtendedKey) (*XprivAuthenticator, error) {
	if xpriv == nil {
		return nil, ErrBip32ExtendedKey
	}
	x := XprivAuthenticator{
		xpriv:    xpriv,
		xpubAuth: &XpubAuthenticator{hdKey: xpriv},
	}
	return &x, nil
}

// NewAccessKeyAuthenticator creates an AccessKeyAuthenticator using an EC private key.
func NewAccessKeyAuthenticator(accessKey *ec.PrivateKey) (*AccessKeyAuthenticator, error) {
	if accessKey == nil {
		return nil, ErrEcPrivateKey
	}
	a := AccessKeyAuthenticator{
		priv: accessKey,
		pub:  accessKey.PubKey(),
	}
	return &a, nil
}

// NewXpubOnlyAuthenticator creates an XpubAuthenticator using a BIP32 extended public key.
func NewXpubOnlyAuthenticator(xpub *bip32.ExtendedKey) (*XpubAuthenticator, error) {
	if xpub == nil {
		return nil, ErrBip32ExtendedKey
	}
	x := XpubAuthenticator{hdKey: xpub}
	return &x, nil
}

var (
	// ErrBip32ExtendedKey indicates that a required BIP32 extended key was not provided.
	ErrBip32ExtendedKey = errors.New("authenticator failed: expected a BIP32 extended key but none was provided")
	// ErrEcPrivateKey indicates that a required EC private key was not provided.
	ErrEcPrivateKey = errors.New("authenticator failed: expected an EC private key but none was provided")
)
