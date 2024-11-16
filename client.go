package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	bip32 "github.com/bitcoin-sv/go-sdk/compat/bip32"
	ec "github.com/bitcoin-sv/go-sdk/primitives/ec"
	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/configs"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/contacts"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/invitations"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/transactions"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/auth"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
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
	contactsAPI     *contacts.API
	invitationsAPI  *invitations.API
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

// Contacts retrieves the full list of user contacts. This method sends an HTTP GET request
// to the "api/v1/contacts" endpoint. The server's response is expected to be unmarshaled into
// a slice of *response.Contact structs. If the request fails or the response cannot be decoded,
// an error is returned.
func (c *Client) Contacts(ctx context.Context) ([]*response.Contact, error) {
	res, err := c.contactsAPI.Contacts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve contacts from user contacts API: %w", err)
	}

	return res, nil
}

// ContactWithPaymail retrieves a specific user contact by their paymail address.
// This method sends an HTTP GET request to "api/v1/contacts/paymail_address", replacing paymail_address
// with the provided paymail. The response is expected to be unmarshaled into a *response.Contact struct.
// If the request fails or the response cannot be decoded, an error is returned.
func (c *Client) ContactWithPaymail(ctx context.Context, paymail string) (*response.Contact, error) {
	res, err := c.contactsAPI.ContactWithPaymail(ctx, paymail)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve contact with paymail from user contacts API: %w", err)
	}

	return res, nil
}

// UpsertContact adds or updates a user contact using the user contacts API.
// This method sends an HTTP PUT request to "api/v1/contacts/paymail_address", replacing paymail_address
// with the user's paymail. The response is expected to be unmarshaled into a *response.Contact struct.
// If the request fails or the response cannot be decoded, an error is returned.
func (c *Client) UpsertContact(ctx context.Context, cmd commands.UpsertContact) (*response.Contact, error) {
	res, err := c.contactsAPI.UpsertContact(ctx, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to upsert contact by calling the user contacts API: %w", err)
	}

	return res, nil
}

// RemoveContact deletes a user contact using the user contacts API.
// This method sends an HTTP DELETE request to "/api/v1/contacts/paymail_address", replacing paymail_address
// with the user's paymail. If the request fails or the response cannot be decoded, an error is returned.
func (c *Client) RemoveContact(ctx context.Context, paymail string) error {
	err := c.contactsAPI.RemoveContact(ctx, paymail)
	if err != nil {
		return fmt.Errorf("failed to remove contact by calling the user contacts API: %w", err)
	}

	return nil
}

// ConfirmContact confirms a user contact using the user contacts API.
// This method sends an HTTP POST request to "/api/v1/contacts/paymail_address", replacing paymail_address
// with the user's paymail. If the request fails or the response cannot be decoded, an error is returned.
func (c *Client) ConfirmContact(ctx context.Context, paymail string) error {
	err := c.contactsAPI.ConfirmContact(ctx, paymail)
	if err != nil {
		return fmt.Errorf("failed to confirm contact by calling the user contacts API: %w", err)
	}

	return nil
}

// UnconfirmContact unconfirms a user contact using the user contacts API.
// This method sends an HTTP DELETE request to "/api/v1/contacts/paymail_address", replacing paymail_address
// with the user's paymail. If the request fails or the response cannot be decoded, an error is returned.
func (c *Client) UnconfirmContact(ctx context.Context, paymail string) error {
	err := c.contactsAPI.UnconfirmContact(ctx, paymail)
	if err != nil {
		return fmt.Errorf("failed to unconfirm contact by calling the user contacts API: %w", err)
	}

	return nil
}

// AcceptInvitation accepts an invitation to add a contact using the user invitations API.
// This method sends an HTTP POST request to "/api/v1/invitations/paymail_address/contacts", replacing paymail_address
// with the user's paymail. If the request fails or the response cannot be decoded, an error is returned.
func (c *Client) AcceptInvitation(ctx context.Context, paymail string) error {
	err := c.invitationsAPI.AcceptInvitation(ctx, paymail)
	if err != nil {
		return fmt.Errorf("failed to accept inivtation by calling the user invitations API: %w", err)
	}

	return nil
}

// RejectInvitation rejects an invitation using the user invitations API.
// This method sends an HTTP DELETE request to "/api/v1/invitations/paymail_address", replacing paymail_address
// with the user's paymail. If the request fails or the response cannot be decoded, an error is returned.
func (c *Client) RejectInvitation(ctx context.Context, paymail string) error {
	err := c.invitationsAPI.RejectInvitation(ctx, paymail)
	if err != nil {
		return fmt.Errorf("failed to reject inivtation by calling the user invitations API: %w", err)
	}

	return nil
}

// SharedConfig retrieves the shared configuration from the user configurations API.
// This method constructs an HTTP GET request to the "api/v1/configs/shared" endpoint and expects
// a response that can be unmarshaled into the response.SharedConfig struct. If the request fails
// or the response cannot be decoded, an error will be returned.
func (c *Client) SharedConfig(ctx context.Context) (*response.SharedConfig, error) {
	res, err := c.configsAPI.SharedConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve shared configuration from user configs API: %w", err)
	}

	return res, nil
}

// DraftTransaction creates a new draft transaction using the user transactions API.
// This method sends an HTTP POST request to the "/draft" endpoint and expects
// a response that can be unmarshaled into a response.DraftTransaction struct.
// If the request fails or the response cannot be decoded, an error is returned.
func (c *Client) DraftTransaction(ctx context.Context, cmd *commands.DraftTransaction) (*response.DraftTransaction, error) {
	res, err := c.transactionsAPI.DraftTransaction(ctx, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to create a draft transaction by calling the user transactions API: %w", err)
	}

	return res, nil
}

// RecordTransaction submits a transaction for recording using the user transactions API.
// This method sends an HTTP POST request to the "/transactions" endpoint, expecting
// a response that can be unmarshaled into a response.Transaction struct.
// If the request fails or the response cannot be decoded, an error is returned.
func (c *Client) RecordTransaction(ctx context.Context, cmd *commands.RecordTransaction) (*response.Transaction, error) {
	res, err := c.transactionsAPI.RecordTransaction(ctx, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to record a transaction with reference ID: %s by calling the user transactions API: %w", cmd.ReferenceID, err)
	}

	return res, nil
}

// UpdateTransactionMetadata updates the metadata of a transaction using the user transactions API.
// This method sends an HTTP PATCH request with updated metadata and expects a response
// that can be unmarshaled into a response.Transaction struct.
// If the request fails or the response cannot be decoded, an error is returned.
func (c *Client) UpdateTransactionMetadata(ctx context.Context, cmd *commands.UpdateTransactionMetadata) (*response.Transaction, error) {
	res, err := c.transactionsAPI.UpdateTransactionMetadata(ctx, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to update a transaction metadata by calling the user user transactions API: %w", err)
	}

	return res, nil
}

// Transactions retrieves a paginated list of transactions from the user transactions API.
// The returned response includes transactions and pagination details, such as the page number,
// sort order, and sorting field (sortBy).
//
// This method allows optional query parameters to be applied via the provided query options.
// The response is expected to unmarshal into a *response.PageModel[response.Transaction] struct.
// If the API request fails or the response cannot be decoded successfully, an error is returned.
func (c *Client) Transactions(ctx context.Context, opts ...queries.TransactionsQueryOption) (*queries.TransactionPage, error) {
	res, err := c.transactionsAPI.Transactions(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve transactions page from the user transactions API: %w", err)
	}

	return res, nil
}

// Transaction retrieves a specific transaction by its ID using the user transactions API.
// This method expects a response that can be unmarshaled into a response.Transaction struct.
// If the request fails or the response cannot be decoded, an error is returned.
func (c *Client) Transaction(ctx context.Context, ID string) (*response.Transaction, error) {
	res, err := c.transactionsAPI.Transaction(ctx, ID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve transaction with ID: %s from the user transactions API: %w", ID, err)
	}

	return res, nil
}

// ErrUnrecognizedAPIResponse indicates that the response received from the SPV Wallet API
// does not match the expected expected format or structure.
var ErrUnrecognizedAPIResponse = errors.New("unrecognized response from API")

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
		contactsAPI:     contacts.NewAPI(cfg.Addr, restyCli),
		invitationsAPI:  invitations.NewAPI(cfg.Addr, restyCli),
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
