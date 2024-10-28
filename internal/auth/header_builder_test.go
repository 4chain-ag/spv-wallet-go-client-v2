package auth_test

import (
	"net/http"
	"testing"

	bip32 "github.com/bitcoin-sv/go-sdk/compat/bip32"
	compat "github.com/bitcoin-sv/go-sdk/compat/bip32"
	ec "github.com/bitcoin-sv/go-sdk/primitives/ec"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/auth"
	"github.com/stretchr/testify/require"
)

func TestHeaderBuilder_New(t *testing.T) {
	tests := map[string]struct {
		expectedErr error
		cfg         *auth.HeaderConfig
	}{
		"build: empty header cfg": {
			expectedErr: auth.ErrMissingHeaderConfig,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			b, err := auth.NewHeaderBuilder(tc.cfg)
			require.Nil(t, b)
			require.ErrorIs(t, tc.expectedErr, err)
		})
	}
}

func TestHeaderBuilder_BuildWithoutBody(t *testing.T) {
	tests := map[string]struct {
		expectedErr    error
		cfg            *auth.HeaderConfig
		requireHeaders HTTPHeadersTestHelper
	}{
		"build: extended key and sign request fields are set": {
			requireHeaders: ExtendedKeyWithSignRequestHeaders,
			cfg: &auth.HeaderConfig{
				ExtendedKey: ExtendedKey(t),
				SignRequest: true,
			},
		},
		"build: extended key set only": {
			requireHeaders: ExtendedKeyWithoutSignRequestHeaders,
			cfg: &auth.HeaderConfig{
				ExtendedKey: ExtendedKey(t),
			},
		},
		"build: private key set only": {
			requireHeaders: PrivateKeyHeaders,
			cfg: &auth.HeaderConfig{
				PrivateKey: PrivateKey(t),
			},
		},
		"build: empty header cfg": {
			expectedErr: auth.ErrMissingKeys,
			cfg:         &auth.HeaderConfig{},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			b, err := auth.NewHeaderBuilder(tc.cfg)
			require.NoError(t, err)

			got, err := b.BuildWithoutBody()
			require.ErrorIs(t, err, tc.expectedErr)
			if tc.requireHeaders != nil {
				tc.requireHeaders(t, got)
			}
		})
	}
}

func TestHeaderBuilder_Build(t *testing.T) {
	tests := map[string]struct {
		expectedErr    error
		body           string
		cfg            *auth.HeaderConfig
		requireHeaders HTTPHeadersTestHelper
	}{
		"build: extended key and sign request fields are set": {
			body:           "{ content: 123 }",
			requireHeaders: ExtendedKeyWithSignRequestHeaders,
			cfg: &auth.HeaderConfig{
				ExtendedKey: ExtendedKey(t),
				SignRequest: true,
			},
		},
		"build: extended key set only": {
			body:           "{ content: 123 }",
			requireHeaders: ExtendedKeyWithoutSignRequestHeaders,
			cfg: &auth.HeaderConfig{
				ExtendedKey: ExtendedKey(t),
			},
		},
		"build: private key set only": {
			body:           "{ content: 123 }",
			requireHeaders: PrivateKeyHeaders,
			cfg: &auth.HeaderConfig{
				PrivateKey: PrivateKey(t),
			},
		},
		"build: empty header cfg": {
			expectedErr: auth.ErrMissingKeys,
			cfg:         &auth.HeaderConfig{},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			b, err := auth.NewHeaderBuilder(tc.cfg)
			require.NoError(t, err)

			got, err := b.Build(tc.body)
			require.ErrorIs(t, err, tc.expectedErr)
			if tc.requireHeaders != nil {
				tc.requireHeaders(t, got)
			}
		})
	}
}

// HTTPHeadersTestHelper is a helper function type for testing HTTP headers.
// It takes a testing instance (`*testing.T`) and an `http.Header` to perform
// validation or assertion tasks on the provided headers.
type HTTPHeadersTestHelper func(t *testing.T, h http.Header)

// ExtendedKey generates a new `*bip32.ExtendedKey` for testing purposes using a
// recommended seed length. This helper function is useful for testing
// extended key-based functionality and cryptographic signing with dynamically
// generated keys.
func ExtendedKey(t *testing.T) *bip32.ExtendedKey {
	t.Helper()
	key, err := bip32.GenerateHDKey(compat.RecommendedSeedLength)
	if err != nil {
		t.Fatalf("test helper - failed to generate HD key from string: %sw", err)
	}
	return key
}

// PrivateKey generates a new `*ec.PrivateKey` for testing purposes.
// This helper simplifies testing of private key-based signing.
// If key creation fails, the test is immediately terminated.
func PrivateKey(t *testing.T) *ec.PrivateKey {
	t.Helper()
	key, err := ec.NewPrivateKey()
	if err != nil {
		t.Fatalf("test helper - failed to create private key: %s", err)
	}
	return key
}

// PrivateKeyHeaders validates the presence of required HTTP headers
// when using a private key for signing. The test will fail if expected
// headers do not match the actual headers.
func PrivateKeyHeaders(t *testing.T, h http.Header) {
	t.Helper()
	expected := []string{
		"X-Auth-Key",
		"X-Auth-Hash",
		"X-Auth-Nonce",
		"X-Auth-Time",
		"X-Auth-Signature",
	}
	got := make([]string, 0, len(expected))
	for k := range h {
		got = append(got, k)
	}
	require.ElementsMatch(t, expected, got)
}

// ExtendedKeyWithSignRequestHeaders validates the presence of required
// HTTP headers when using an extended key for signing with signing enabled.
// The test will fail if expected headers do not match the actual headers.
func ExtendedKeyWithSignRequestHeaders(t *testing.T, h http.Header) {
	t.Helper()
	expected := []string{
		"X-Auth-Xpub",
		"X-Auth-Hash",
		"X-Auth-Nonce",
		"X-Auth-Time",
		"X-Auth-Signature",
	}
	got := make([]string, 0, len(expected))
	for k := range h {
		got = append(got, k)
	}
	require.ElementsMatch(t, expected, got)
}

// ExtendedKeyWithoutSignRequestHeaders validates the presence of required
// HTTP headers when using an extended key for signing with signing disabled.
// The test will fail if expected headers do not match the actual headers.
func ExtendedKeyWithoutSignRequestHeaders(t *testing.T, h http.Header) {
	t.Helper()
	expected := []string{"X-Auth-Xpub"}
	got := make([]string, 0, len(expected))
	for k := range h {
		got = append(got, k)
	}
	require.ElementsMatch(t, expected, got)
}
