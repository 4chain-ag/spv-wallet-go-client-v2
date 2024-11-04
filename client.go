package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	bip32 "github.com/bitcoin-sv/go-sdk/compat/bip32"
	ec "github.com/bitcoin-sv/go-sdk/primitives/ec"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/configs"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/transactions"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/auth"
	"github.com/bitcoin-sv/spv-wallet-go-client/query"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/go-resty/resty/v2"
)

// Config holds configuration settings for establishing a connection and handling
// request details in the application.
type Config struct {
	Addr      string            // The base address of the SPV Wallet API.
	Timeout   time.Duration     // The HTTP requests timeout duration.
	Transport http.RoundTripper // Custom HTTP transport, allowing optional customization of the HTTP client behavior.
}

// NewDefaultConfig returns a default configuration for connecting to the SPV Wallet API,
// setting a one-minute timeout, using the default HTTP transport, and applying the
// base API address as the addr value.
func NewDefaultConfig(addr string) Config {
	return Config{
		Addr:      addr,
		Timeout:   1 * time.Minute,
		Transport: http.DefaultTransport,
	}
}

// Client provides methods for user-related and admin-related APIs.
// This struct is designed to abstract and simplify the process of making HTTP calls
// to the relevant endpoints. By utilizing this Client struct, developers can easily
// interact with both user and admin APIs without needing to manage the details
// of the HTTP requests and responses directly.
type Client struct {
	configsAPI      *configs.API
	transactionsAPI *transactions.API
}

// NewWithXPub creates a new client instance using an extended public key (xPub).
// Requests made with this instance will not be signed, that's why we strongly recommend to use `WithXPriv` or `WithAccessKey` option instead.
func NewWithXPub(cfg Config, xPub string) (*Client, error) {
	key, err := bip32.GetHDKeyFromExtendedPublicKey(xPub)
	if err != nil {
		return nil, fmt.Errorf("failed to generate HD key from xPub: %w", err)
	}

	authenticator, err := auth.NewXpubOnlyAuthenticator(key)
	if err != nil {
		return nil, fmt.Errorf("failed to intialized xpub authenticator: %w", err)
	}

	return newClient(cfg, authenticator), nil
}

// NewWithXPriv creates a new client instance using an extended private key (xPriv).
// Generates an HD key from the provided xPriv and sets up the client instance to sign requests
// by setting the SignRequest flag to true. The generated HD key can be used for secure communications.
func NewWithXPriv(cfg Config, xPriv string) (*Client, error) {
	key, err := bip32.GenerateHDKeyFromString(xPriv)
	if err != nil {
		return nil, fmt.Errorf("failed to generate HD key from xpriv: %w", err)
	}

	authenticator, err := auth.NewXprivAuthenticator(key)
	if err != nil {
		return nil, fmt.Errorf("failed to intialized xpriv authenticator: %w", err)
	}

	return newClient(cfg, authenticator), nil
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

	authenticator, err := auth.NewAccessKeyAuthenticator(key)
	if err != nil {
		return nil, fmt.Errorf("failed to intialized access key authenticator: %w", err)
	}

	return newClient(cfg, authenticator), nil
}

// SharedConfig retrieves the shared configuration from the user configurations API.
// This method constructs an HTTP GET request to the "/shared" endpoint and expects
// a response that can be unmarshaled into the response.SharedConfig struct.
// If the request fails or the response cannot be decoded, an error will be returned.
func (c *Client) SharedConfig(ctx context.Context) (*response.SharedConfig, error) {
	res, err := c.configsAPI.SharedConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve shared configuration from user configs API: %w", err)
	}

	return res, nil
}

// RecordTransactionArgs holds the arguments required to record a user transaction.
// It contains metadata, the hex representation of the transaction, and the reference ID.
type RecordTransactionArgs struct {
	Metadata    query.Metadata // Metadata associated with the transaction.
	Hex         string         // Hexadecimal string representation of the transaction.
	ReferenceID string         // Reference ID for the transaction.
}

// ParseToRecordTransactionRequest converts RecordTransactionArgs to a RecordTransactionRequest
// for SPV Wallet API consumption.
func (r RecordTransactionArgs) ParseToRecordTransactionRequest() transactions.RecordTransactionRequest {
	return transactions.RecordTransactionRequest{
		Metadata:    r.Metadata,
		Hex:         r.Hex,
		ReferenceID: r.ReferenceID,
	}
}

// DraftTransactionArgs holds the arguments required to create user draft transaction.
// It includes the transaction configuration and associated metadata.
type DraftTransactionArgs struct {
	Config   response.TransactionConfig // Configuration for the transaction.
	Metadata query.Metadata             // Metadata related to the transaction.
}

// ParseToDraftTransactionRequest converts DraftTransactionArgs to a DraftTransactionRequest
// for SPV Wallet API consumption.
func (d DraftTransactionArgs) ParseToDraftTransactionRequest() transactions.DraftTransactionRequest {
	return transactions.DraftTransactionRequest{
		Config:   d.Config,
		Metadata: d.Metadata,
	}
}

// UpdateTransactionMetadataArgs holds the arguments required to update a user transaction's metadata.
// It contains the transaction ID and the new metadata.
type UpdateTransactionMetadataArgs struct {
	ID       string         // Unique identifier of the transaction to be updated.
	Metadata query.Metadata // New metadata to associate with the transaction.
}

// ParseUpdateTransactionMetadataRequest converts UpdateTransactionMetadataArgs to an
// UpdateTransactionMetadataRequest for SPV Wallet API consumption.
func (u UpdateTransactionMetadataArgs) ParseUpdateTransactionMetadataRequest() transactions.UpdateTransactionMetadataRequest {
	return transactions.UpdateTransactionMetadataRequest{
		ID:       u.ID,
		Metadata: u.Metadata,
	}
}

func (c *Client) DraftTransaction(ctx context.Context, args DraftTransactionArgs) (*response.DraftTransaction, error) {
	res, err := c.transactionsAPI.DraftTransaction(ctx, args.ParseToDraftTransactionRequest())
	if err != nil {
		return nil, fmt.Errorf("failed to create draft transaction by call user transactions API: %w", err)
	}

	return res, nil
}

func (c *Client) RecordTransaction(ctx context.Context, args RecordTransactionArgs) (*response.Transaction, error) {
	res, err := c.transactionsAPI.RecordTransaction(ctx, args.ParseToRecordTransactionRequest())
	if err != nil {
		return nil, fmt.Errorf("failed to record transaction with reference ID: %s by call user transactions API: %w", args.ReferenceID, err)
	}

	return res, nil
}

func (c *Client) UpdateTransactionMetadata(ctx context.Context, args UpdateTransactionMetadataArgs) (*response.Transaction, error) {
	res, err := c.transactionsAPI.UpdateTransactionMetadata(ctx, args.ParseUpdateTransactionMetadataRequest())
	if err != nil {
		return nil, fmt.Errorf("failed to update transactions metadata by call user user transactions API: %w", err)
	}

	return res, nil
}

func (c *Client) Transactions(ctx context.Context, opts ...query.BuilderOption) ([]*response.Transaction, error) {
	res, err := c.transactionsAPI.Transactions(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve transactions from user transactions API: %w", err)
	}

	return res, nil
}

func (c *Client) Transaction(ctx context.Context, ID string) (*response.Transaction, error) {
	res, err := c.transactionsAPI.Transaction(ctx, ID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve transaction with ID: %s from user transactions API: %w", ID, err)
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
		return nil, errors.Join(err1, err2)
	}

	return pk, nil
}

type authenticator interface {
	Authenticate(r *resty.Request) error
}

func newClient(cfg Config, auth authenticator) *Client {
	restyCli := newRestyClient(cfg, auth)
	cli := Client{
		configsAPI:      configs.NewAPI(cfg.Addr, restyCli),
		transactionsAPI: transactions.NewAPI(cfg.Addr, restyCli),
	}
	return &cli
}

func newRestyClient(cfg Config, auth authenticator) *resty.Client {
	return resty.New().
		SetTransport(cfg.Transport).
		SetBaseURL(cfg.Addr).
		SetTimeout(cfg.Timeout).
		OnBeforeRequest(func(_ *resty.Client, r *resty.Request) error {
			return auth.Authenticate(r)
		}).
		SetError(&models.SPVError{}).
		OnAfterResponse(func(_ *resty.Client, r *resty.Response) error {
			if r.IsSuccess() {
				return nil
			}

			if spvError, ok := r.Error().(*models.SPVError); ok && len(spvError.Code) > 0 {
				return spvError
			}

			return fmt.Errorf("%w: %s", ErrUnrecognizedAPIResponse, r.Body())
		})
}

// ErrUnrecognizedAPIResponse indicates that the response received from the SPV Wallet API
// does not match the expected expected format or structure.
var ErrUnrecognizedAPIResponse = errors.New("unrecognized response from API")
