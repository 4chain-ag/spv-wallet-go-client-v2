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

// Actor represents an individual user or entity within the SPV Wallet ecosystem.
// It includes details like the alias, private key (xPriv), public key (xPub), and paymail address.
type Actor struct {
	alias   string // The unique alias for the actor.
	xPriv   string // The extended private key for the actor.
	xPub    string // The extended public key for the actor.
	paymail string // The paymail address associated with the actor.
}

// NewActor creates and returns a new Actor instance with a unique alias and domain.
// It generates a new set of random wallet keys (xPriv and xPub) and constructs the paymail address.
// On success, it returns the initialized Actor instance and a nil error.
// If key generation fails, it returns a non-nil error.
func NewActor(alias, domain string) (*Actor, error) {
	keys, err := walletkeys.RandomKeys()
	if err != nil {
		return nil, fmt.Errorf("failed to generate wallet keys: %w", err)
	}
	return &Actor{
		alias:   alias,
		xPriv:   keys.XPriv(),
		xPub:    keys.XPub(),
		paymail: alias + "@" + domain,
	}, nil
}

// ActorService provides an interface to interact with the SPV Wallet API.
// It supports operations such as retrieving the balance, accessing the paymail,
// fetching transactions, and transferring funds within the wallet ecosystem.
type ActorService struct {
	userAPI *wallet.UserAPI // The API client for interacting with the SPV Wallet.
	actor   *Actor          // The actor associated with this service.
}

// Paymail returns the paymail address of the associated actor.
func (a *ActorService) Paymail() string {
	return a.actor.paymail
}

// Balance retrieves the current satoshi balance for the actor.
// It accepts a context parameter to manage cancellation and timeouts.
// On success, it returns the current balance and a nil error.
// If the API call fails, it returns a non-nil error with details of the failure.
func (a *ActorService) Balance(ctx context.Context) (uint64, error) {
	xpub, err := a.userAPI.XPub(ctx)
	if err != nil {
		return 0, fmt.Errorf("could not fetch xPub to retrieve current balance: %w", err)
	}
	return xpub.CurrentBalance, nil
}

// Transactions retrieves the list of transactions associated with the actor's wallet.
// It accepts a context parameter to manage cancellation and timeouts.
// On success, it returns a slice of transactions and a nil error.
// If the API call fails, it returns a non-nil error with details of the failure.
func (a *ActorService) Transactions(ctx context.Context) ([]*response.Transaction, error) {
	page, err := a.userAPI.Transactions(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not fetch transactions: %w", err)
	}
	return page.Content, nil
}

// TransferFunds sends a specified amount of satoshis to a recipient's paymail.
// It accepts a context parameter to manage cancellation and timeouts.
// On success, it returns the transaction representing the fund transfer and a nil error.
// If the actor has insufficient funds or the API call fails, it returns a non-nil error.
func (a *ActorService) TransferFunds(ctx context.Context, paymail string, funds uint64) (*response.Transaction, error) {
	balance, err := a.Balance(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not fetch balance: %w", err)
	}
	if balance < funds {
		return nil, fmt.Errorf("insufficient balance: %d to spend: %d in transaction", balance, funds)
	}

	recipient := commands.Recipients{To: paymail, Satoshis: funds}
	transaction, err := a.userAPI.SendToRecipients(ctx, &commands.SendToRecipients{
		Recipients: []*commands.Recipients{&recipient},
		Metadata:   map[string]any{"description": "regression-test"},
	})
	if err != nil {
		return nil, fmt.Errorf("could not transfer funds: %w", err)
	}
	return transaction, nil
}

// NewActorService initializes and returns a new ActorService instance.
// It accepts the API URL and an Actor instance as input.
// On success, it returns the initialized ActorService and a nil error.
// If the user API initialization fails, it returns a non-nil error with details of the failure.
func NewActorService(url string, actor *Actor) (*ActorService, error) {
	userAPI, err := wallet.NewUserAPIWithXPriv(config.New(config.WithAddr(url)), actor.xPriv)
	if err != nil {
		return nil, fmt.Errorf("could not initialize user API: %w", err)
	}
	return &ActorService{userAPI: userAPI, actor: actor}, nil
}
