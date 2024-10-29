package client

import (
	"context"
	"errors"
	"fmt"
	"time"

	bip32 "github.com/bitcoin-sv/go-sdk/compat/bip32"
	ec "github.com/bitcoin-sv/go-sdk/primitives/ec"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/auth"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/httpx"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

// DefaultConfig defines default client options for local development.
// This configuration is recommended for testing and development environments
// only, as it uses a localhost address and a shorter timeout duration that
// may not be suitable for production.
var DefaultConfig = Config{
	Addr:    "http://localhost:3003",
	Timeout: 1 * time.Minute,
}

// Config holds configuration settings for establishing a connection and handling
// request details in the application.
type Config struct {
	Addr    string        // The base address of the SPV Wallet API.
	Timeout time.Duration // Timeout duration for connection attempts.
}

// Client provides methods for user-related and admin-related APIs.
// This struct is designed to abstract and simplify the process of making HTTP calls
// to the relevant endpoints. By utilizing this Client struct, developers can easily
// interact with both user and admin APIs without needing to manage the details
// of the HTTP requests and responses directly.
type Client struct {
	userAPI *user.API
}

// SharedConfig retrieves the shared configuration from the user configurations API.
// This method constructs an HTTP GET request to the "/shared" endpoint and expects
// a response that can be unmarshaled into the response.SharedConfig struct.
// If the request fails or the response cannot be decoded, an error will be returned.
func (c *Client) SharedConfig(ctx context.Context) (*response.SharedConfig, error) {
	res, err := c.userAPI.ConfigsAPI.SharedConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve shared configuration from User Configs API: %w", err)
	}
	return res, nil
}

func privateKeyFromHexOrWIF(s string) (*ec.PrivateKey, error) {
	pk, err1 := ec.PrivateKeyFromWif(s)
	if err1 == nil {
		return pk, nil
	}
	pk, err2 := ec.PrivateKeyFromHex(s)
	if err2 != nil {
		return nil, errors.Join(err1, err2) // Join the errors from both attempts.
	}
	return pk, nil
}

// newClient initializes a new client instance with the given address and configuration.
// This function creates a Resty HTTP client with the specified address and configuration
// and initializes the client instance allowing for interaction with user-related and admin-related endpoints.
func newClient(cfg httpx.Config) (*Client, error) {
	cli, err := httpx.NewResty(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create Resty client: %w", err)
	}
	return &Client{userAPI: user.NewAPI(cfg.Addr, cli)}, nil
}

// NewWithXPub creates a new client instance using an extended public key (xPub).
// Generates a hierarchical deterministic (HD) key from the provided xPub and constructs
// the necessary configuration for the API instance. The SignRequest flag is set to false,
// indicating that requests made with this instance will not be signed.
func NewWithXPub(cfg Config, xPub string) (*Client, error) {
	key, err := bip32.GetHDKeyFromExtendedPublicKey(xPub)
	if err != nil {
		return nil, fmt.Errorf("failed to generate HD key from xPub: %w", err)
	}
	return newClient(httpx.Config{
		Addr:    cfg.Addr,
		Timeout: cfg.Timeout,
		HeaderConfig: auth.HeaderConfig{
			ExtendedKey: key,
			SignRequest: false,
		},
	})
}

// NewWithXPriv creates a new client instance using an extended private key (xPriv).
// Generates an HD key from the provided xPriv and sets up the client instance to sign requests
// by setting the SignRequest flag to true. The generated HD key can be used for secure communications.
func NewWithXPriv(cfg Config, xPriv string) (*Client, error) {
	key, err := bip32.GenerateHDKeyFromString(xPriv)
	if err != nil {
		return nil, fmt.Errorf("failed to generate HD key from xPriv: %w", err)
	}
	return newClient(httpx.Config{
		Addr:    cfg.Addr,
		Timeout: cfg.Timeout,
		HeaderConfig: auth.HeaderConfig{
			ExtendedKey: key,
			SignRequest: true,
		},
	})
}

// NewWithAccessKey creates a new client instance using an access key.
// Function attempts to convert the provided access key from either hex or WIF format
// to a PrivateKey. The resulting PrivateKey is used to sign requests made by the client instance
// by setting the SignRequest flag to true.
func NewWithAccessKey(cfg Config, accessKey string) (*Client, error) {
	key, err := privateKeyFromHexOrWIF(accessKey)
	if err != nil {
		return nil, fmt.Errorf("failed to return private key from hex or WIF: %w", err)
	}
	return newClient(httpx.Config{
		Addr:    cfg.Addr,
		Timeout: cfg.Timeout,
		HeaderConfig: auth.HeaderConfig{
			PrivateKey:  key,
			SignRequest: true,
		},
	})
}
