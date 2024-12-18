package regressiontests

import (
	"context"
	"fmt"

	wallet "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/config"
	"github.com/bitcoin-sv/spv-wallet-go-client/walletkeys"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

// transactionsSlice represents a slice of response.Transaction objects.
type transactionsSlice []*response.Transaction

// Has checks if a transaction with the specified ID exists in the transactions slice.
// It returns true if a transaction with the given ID is found, and false otherwise.
func (tt transactionsSlice) Has(id string) bool {
	for _, t := range tt {
		if t.ID == id {
			return true
		}
	}

	return false
}

// user represents an individual user within the SPV Wallet ecosystem.
// It includes details like the alias, private key (xPriv), public key (xPub), and paymail address.
// The user struct also utilizes the wallet's UserAPI client to interact with the SPV Wallet API
// for transaction-related operations.
type user struct {
	alias   string          // The unique alias for the actor.
	xPriv   string          // The extended private key for the actor.
	xPub    string          // The extended public key for the actor.
	paymail string          // The paymail address associated with the actor.
	client  *wallet.UserAPI // The API client for interacting with the SPV Wallet.
}

// transferFunds sends a specified amount of satoshis to a recipient's paymail.
// It accepts a context parameter to manage cancellation and timeouts.
// On success, it returns the transaction representing the fund transfer and a nil error.
// If the actor has insufficient funds or the API call fails, it returns a non-nil error.
func (u *user) transferFunds(ctx context.Context, paymail string, funds uint64) (*response.Transaction, error) {
	balance, err := u.balance(ctx)
	if err != nil {
		return nil, fmt.Errorf("balance failed: %w", err)
	}
	if balance < funds {
		return nil, fmt.Errorf("insufficient balance: %d available, %d required", balance, funds)
	}

	recipient := commands.Recipients{To: paymail, Satoshis: funds}
	transaction, err := u.client.SendToRecipients(ctx, &commands.SendToRecipients{
		Recipients: []*commands.Recipients{&recipient},
		Metadata:   map[string]any{"description": "regression-test"},
	})
	if err != nil {
		return nil, fmt.Errorf("could not transfer funds to %s: %w", paymail, err)
	}

	return transaction, nil
}

// transactions retrieves the list of transactions associated with the given user.
// It accepts a context parameter to manage cancellation and timeouts.
// On success, it returns a slice of transactions and a nil error.
// If the API call fails, it returns a non-nil error with details of the failure.
func (u *user) transactions(ctx context.Context) (transactionsSlice, error) {
	page, err := u.client.Transactions(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not retrive transactions: %w", err)
	}

	return page.Content, nil
}

// balance retrieves the current satoshi balance for given actor.
// It accepts a context parameter to manage cancellation and timeouts.
// On success, it returns the current balance and a nil error.
// If the API call fails, it returns a non-nil error with details of the failure.
func (u *user) balance(ctx context.Context) (uint64, error) {
	xPub, err := u.client.XPub(ctx)
	if err != nil {
		return 0, fmt.Errorf("could not retrive xPub: %w", err)
	}

	return xPub.CurrentBalance, nil
}

// initUser initializes a new user within the SPV Wallet ecosystem.
// It accepts the alias and SPV Wallet API URL as input parameters.
// The function generates a random pair of wallet keys (xPub, xPriv) and uses the xPriv key
// to initialize the wallet's client, enabling transaction-related operations.
// On success, it returns the initialized user and a nil error.
// If user initialization fails, it returns a non-nil error with details of the failure.
func initUser(alias, url string) (*user, error) {
	keys, err := walletkeys.RandomKeys()
	if err != nil {
		return nil, fmt.Errorf("could not generate random keys: %w", err)
	}

	client, err := wallet.NewUserAPIWithXPriv(config.New(config.WithAddr(url)), keys.XPriv())
	if err != nil {
		return nil, fmt.Errorf("could not initialize user API for alias %q: %w", alias, err)
	}

	domain, err := fetchPaymailDomain(client)
	if err != nil {
		return nil, fmt.Errorf("could not initialize user %q: %w", alias, err)
	}

	return &user{
		alias:   alias,
		xPriv:   keys.XPriv(),
		xPub:    keys.XPub(),
		paymail: fmt.Sprintf("%s@%s", alias, domain),
		client:  client,
	}, nil
}

// initLeaderUser initializes a new user representing the "Leader" role in the SPV Wallet ecosystem.
// It accepts the SPV Wallet API URL and the xPriv key for the Leader account.
// The function initializes a UserAPI client using the provided xPriv key, enabling transaction-related interactions.
// On success, it returns the initialized user with the alias "Leader" and a nil error.
// If the initialization fails, a non-nil error with details of the failure is returned.
func initLeaderUser(url, xPriv string) (*user, error) {
	client, err := wallet.NewUserAPIWithXPriv(config.New(config.WithAddr(url)), xPriv)
	if err != nil {
		return nil, fmt.Errorf("could not initialize user API for leader: %w", err)
	}

	domain, err := fetchPaymailDomain(client)
	if err != nil {
		return nil, fmt.Errorf("could not initialize leader user: %w", err)
	}

	return &user{
		alias:   "Leader",
		xPriv:   xPriv,
		paymail: fmt.Sprintf("Leader@%s", domain),
		client:  client,
	}, nil
}

// fetchPaymailDomain retrieves and validates the paymail domain from the wallet's shared configuration.
// It returns the domain as a string if exactly one domain is configured.
// If no domain or multiple domains are configured, it returns an error.
func fetchPaymailDomain(client *wallet.UserAPI) (string, error) {
	sharedCfg, err := client.SharedConfig(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to retrieve shared configuration: %w", err)
	}

	if len(sharedCfg.PaymailDomains) == 0 {
		return "", fmt.Errorf("no paymail domains found in shared configuration")
	}

	if len(sharedCfg.PaymailDomains) > 1 {
		return "", fmt.Errorf("expected one paymail domain, found %d", len(sharedCfg.PaymailDomains))
	}

	return sharedCfg.PaymailDomains[0], nil
}
