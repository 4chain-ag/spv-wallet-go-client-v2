package auth

import (
	"encoding/hex"
	"fmt"
	"net/http"

	bip32 "github.com/bitcoin-sv/go-sdk/compat/bip32"
	ec "github.com/bitcoin-sv/go-sdk/primitives/ec"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/go-resty/resty/v2"
)

type Authenticator interface {
	Authenticate(r *resty.Request) error
}

func NewXpubOnlyAuthenticator(xpub *bip32.ExtendedKey) Authenticator {
	return &xpubAuth{hdKey: xpub}
}

func NewXprivAuthenticator(xpriv *bip32.ExtendedKey) Authenticator {
	return &xprivAuth{
		xpubAuth: &xpubAuth{hdKey: xpriv},
	}
}

func NewAccessKeyAuthenticator(accessKey *ec.PrivateKey) Authenticator {
	return &accessKeyAuth{
		priv: accessKey,
		pub:  accessKey.PubKey(),
	}
}

type xpubAuth struct {
	hdKey *bip32.ExtendedKey
}

func (a *xprivAuth) Authenticate(r *resty.Request) error {
	err := a.xpubAuth.Authenticate(r)
	if err != nil {
		return fmt.Errorf("failed to set xpub header: %w", err)
	}

	body := bodyString(r)
	header := make(http.Header)
	err = setSignature(&header, a.xpriv, body)
	if err != nil {
		return fmt.Errorf("failed to sign request with xpriv: %w", err)
	}
	r.SetHeaderMultiValues(header)
	return nil
}

type xprivAuth struct {
	xpubAuth *xpubAuth
	xpriv    *bip32.ExtendedKey
}

func (a *xpubAuth) Authenticate(r *resty.Request) error {
	xPub, err := bip32.GetExtendedPublicKey(a.hdKey)
	if err != nil {
		return fmt.Errorf("failed to get extended public key: %w", err)
	}
	r.SetHeader(models.AuthHeader, xPub)
	return nil
}

type accessKeyAuth struct {
	priv *ec.PrivateKey
	pub  *ec.PublicKey
}

func (a *accessKeyAuth) Authenticate(r *resty.Request) error {
	header := make(http.Header)
	header.Set(models.AuthAccessKey, a.pubKeyHex())

	body := bodyString(r)

	sign, err := createSignatureAccessKey(a.privKeyHex(), body)
	if err != nil {
		return fmt.Errorf("failed to sign request with access key: %w", err)
	}
	setSignatureHeaders(&header, sign)
	return nil
}

func (a *accessKeyAuth) privKeyHex() string {
	return hex.EncodeToString(a.priv.Serialize())
}

func (a *accessKeyAuth) pubKeyHex() string {
	return hex.EncodeToString(a.pub.SerializeCompressed())
}

func bodyString(r *resty.Request) string {
	switch r.Method {
	case http.MethodGet:
		return ""
	}
	return ""
}
