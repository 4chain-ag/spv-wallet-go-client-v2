package regressiontests

import (
	"context"
	"fmt"

	wallet "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/config"
	"github.com/bitcoin-sv/spv-wallet-go-client/walletkeys"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

// actor represents an individual user or entity within the SPV Wallet ecosystem.
// It includes details like the alias, private key (xPriv), public key (xPub), and paymail address.
type actor struct {
	alias   string // The unique alias for the actor.
	xPriv   string // The extended private key for the actor.
	xPub    string // The extended public key for the actor.
	paymail string // The paymail address associated with the actor.
}

// newActor creates and returns a new Actor instance with a unique alias and domain.
// It generates a new set of random wallet keys (xPriv and xPub) and constructs the paymail address.
// On success, it returns the initialized Actor instance and a nil error.
// If key generation fails, it returns a non-nil error.
func newActor(alias, domain string) (*actor, error) {
	keys, err := walletkeys.RandomKeys()
	if err != nil {
		return nil, fmt.Errorf("failed to generate wallet keys: %w", err)
	}
	return &actor{
		alias:   alias,
		xPriv:   keys.XPriv(),
		xPub:    keys.XPub(),
		paymail: alias + "@" + domain,
	}, nil
}

// actorService provides an interface to interact with the SPV Wallet API.
// It supports operations such as retrieving the balance, accessing the paymail,
// fetching transactions, and transferring funds within the wallet ecosystem.
type actorService struct {
	userAPI         *wallet.UserAPI  // The API client for interacting with the SPV Wallet.
	actor           *actor           // The actor associated with this service.
	transferService *transferService // The API for transferring funds with the SPV Wallet.
}

// paymail returns the paymail address of the associated actor.
func (a *actorService) paymail() string {
	return a.actor.paymail
}

// balance retrieves the current satoshi balance for the actor.
// It accepts a context parameter to manage cancellation and timeouts.
// On success, it returns the current balance and a nil error.
// If the API call fails, it returns a non-nil error with details of the failure.
func (a *actorService) balance(ctx context.Context) (uint64, error) {
	return a.transferService.balance(ctx)
}

// transactions retrieves the list of transactions associated with the actor's wallet.
// It accepts a context parameter to manage cancellation and timeouts.
// On success, it returns a slice of transactions and a nil error.
// If the API call fails, it returns a non-nil error with details of the failure.
func (a *actorService) transactions(ctx context.Context) (transactionsSlice, error) {
	page, err := a.userAPI.Transactions(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not fetch transactions: %w", err)
	}
	return page.Content, nil
}

// transferFunds sends a specified amount of satoshis to a recipient's paymail.
// It accepts a context parameter to manage cancellation and timeouts.
// On success, it returns the transaction representing the fund transfer and a nil error.
// If the actor has insufficient funds or the API call fails, it returns a non-nil error.
func (a *actorService) transferFunds(ctx context.Context, paymail string, funds uint64) (*response.Transaction, error) {
	return a.transferService.transferFunds(ctx, paymail, funds)
}

// newActorService initializes and returns a new ActorService instance.
// It accepts the API URL and an Actor instance as input.
// On success, it returns the initialized ActorService and a nil error.
// If the user API initialization fails, it returns a non-nil error with details of the failure.
func newActorService(url string, actor *actor) (*actorService, error) {
	userAPI, err := wallet.NewUserAPIWithXPriv(config.New(config.WithAddr(url)), actor.xPriv)
	if err != nil {
		return nil, fmt.Errorf("could not initialize user API: %w", err)
	}
	transferService, err := newTransferService(userAPI)
	if err != nil {
		return nil, fmt.Errorf("could not initialize transfer service: %w", err)
	}
	return &actorService{userAPI: userAPI, actor: actor, transferService: transferService}, nil
}
