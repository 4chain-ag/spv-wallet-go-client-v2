package auth

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"

	bip32 "github.com/bitcoin-sv/go-sdk/compat/bip32"
	ec "github.com/bitcoin-sv/go-sdk/primitives/ec"
	"github.com/bitcoin-sv/spv-wallet/models"
)

// HeaderConfig contains configuration settings for request header authentication.
// It defines options for signing requests and selecting the appropriate key for signing.
type HeaderConfig struct {
	SignRequest bool               // SignRequest indicates whether to sign the request.
	PrivateKey  *ec.PrivateKey     // PrivateKey is used for signing when no ExtendedKey is provided.
	ExtendedKey *bip32.ExtendedKey // ExtendedKey is used for signing when provided.
}

// IsExtendedKey checks if an extended key is available for signing.
func (h HeaderConfig) IsExtendedKey() bool { return h.ExtendedKey != nil }

// IsPrivateKey checks if a private key is available for signing.
func (h HeaderConfig) IsPrivateKey() bool { return h.PrivateKey != nil }

// IsSignRequest indicates whether the request should be signed based on
// the configuration.
func (h HeaderConfig) IsSignRequest() bool { return h.SignRequest }

// HeaderBuilder constructs HTTP headers for signing requests.
// It supports building headers using either an ExtendedKey (bip32.ExtendedKey)
// or a PrivateKey (ec.PrivateKey). The header can include a signature
// based on the specified Body and the SignRequest flag.
type HeaderBuilder struct {
	cfg *HeaderConfig
}

// Build generates the HTTP headers based on the keys provided.
// It returns an error if neither an ExtendedKey nor a PrivateKey is set.
// Body is the content to be signed or included in the header.
func (b *HeaderBuilder) Build(body string) (http.Header, error) {
	if b.cfg.IsExtendedKey() {
		return b.buildWithExtendedKey(body)
	}
	if b.cfg.IsPrivateKey() {
		return b.buildWithPrivateKey(body)
	}
	return nil, ErrMissingKeys
}

// BuildWithoutBody constructs an HTTP header without any request body content.
// This method calls the Build method with an empty string as the body, creating
// a header suitable for GET, DELETE, or other HTTP requests that do not require
// a request payload.
func (b *HeaderBuilder) BuildWithoutBody() (http.Header, error) {
	return b.Build("")
}

// buildWithExtendedKey creates headers using the ExtendedKey.
// It returns an error if the ExtendedKey is not set or if
// there is an issue generating the public key or signature.
func (b *HeaderBuilder) buildWithExtendedKey(body string) (http.Header, error) {
	header := make(http.Header)
	xPub, err := bip32.GetExtendedPublicKey(b.cfg.ExtendedKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get extended public key: %w", err)
	}
	header.Set(models.AuthHeader, xPub)

	if b.cfg.IsSignRequest() {
		if err := setSignature(&header, b.cfg.ExtendedKey, body); err != nil {
			return nil, fmt.Errorf("failed to set signature: %w", err)
		}
	}
	return header, nil
}

// buildWithPrivateKey creates headers using the PrivateKey.
// It returns an error if the PrivateKey is not set or if
// there is an issue generating the signature for the access key.
func (b *HeaderBuilder) buildWithPrivateKey(body string) (http.Header, error) {
	header := make(http.Header)
	hex := hex.EncodeToString(b.cfg.PrivateKey.Serialize())
	sign, err := createSignatureAccessKey(hex, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create signature for access key: %w", err)
	}

	header.Set(models.AuthAccessKey, sign.AccessKey)
	setSignatureHeaders(&header, sign)
	return header, nil
}

// NewHeaderBuilder constructs a new HeaderBuilder instance using the provided configuration.
// This function initializes a HeaderBuilder with the given HeaderConfig. The HeaderBuilder
// is responsible for constructing HTTP headers necessary for signing requests. If the
// provided configuration is nil, an error will be returned to indicate that the header
// configuration is required.
func NewHeaderBuilder(cfg *HeaderConfig) (*HeaderBuilder, error) {
	if cfg == nil {
		return nil, ErrMissingHeaderConfig
	}
	return &HeaderBuilder{cfg: cfg}, nil
}

var (
	// ErrMissingKeys is returned when the HeaderBuilder does not have either an ExtendedKey (bip32.ExtendedKey)
	// or a PrivateKey (ec.PrivateKey) set, which are required for building the HTTP auth headers.
	ErrMissingKeys = errors.New("Header builder requires either an ExtendedKey (bip32.ExtendedKey) or PrivateKey (ec.PrivateKey) to build auth headers")
	// ErrMissingHeaderConfig is returned when NewHeaderBuilder is called with a nil `HeaderConfig` argument,
	// indicating that header authentication configuration is required but not provided.
	ErrMissingHeaderConfig = errors.New("Header builder requires header auth config to build HTTP auth headers")
)
