package regressiontests

import (
	"context"
	"errors"
	"fmt"

	wallet "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/config"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

// leaderServiceConfig contains configuration settings for initializing a LeaderService.
// These include the environment URL and private keys required for admin and user operations.
type leaderServiceConfig struct {
	envURL     string // URL of the SPV Wallet API environment.
	envXPriv   string // Extended private key (xPriv) for the user account.
	adminXPriv string // Extended private key (xPriv) for the admin account.
}

// leaderService provides an interface to interact with the SPV Wallet API.
// It supports operations like retrieving balances, managing transactions,
// creating and managing actors, and transferring funds within the wallet ecosystem.
type leaderService struct {
	cfg             *leaderServiceConfig // Configuration for the service.
	transferService *transferService     // The API for transferring funds with the SPV Wallet.
	adminAPI        *wallet.AdminAPI     // Admin API client for privileged operations.
	userAPI         *wallet.UserAPI      // User API client for user-level operations.
	actors          []*actor             // List of actors created by this service.
	domain          string               // Paymail domain used by the service.
}

// paymail returns the paymail address of the leader.
// The paymail is constructed as "Leader@" followed by the configured domain.
func (l *leaderService) paymail() string {
	return "Leader@" + l.domain
}

// balance retrieves the current satoshi balance for the leader account.
// It accepts a context parameter to manage cancellation and timeouts.
// The method fetches the xPub key of the user account and uses it to determine the current balance.
// On success, it returns the balance in satoshis and a nil error. If the API call fails, it returns a non-nil error.
func (l *leaderService) balance(ctx context.Context) (uint64, error) {
	return l.transferService.balance(ctx)
}

// transactions retrieves the list of transactions associated with the leader's account.
// It accepts a context parameter to manage cancellation and timeouts.
// On success, it returns a slice of transactions contained in the paginated response and a nil error.
// If the API call fails, it returns a non-nil error.
func (l *leaderService) transactions(ctx context.Context) (transactionsSlice, error) {
	page, err := l.userAPI.Transactions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transactions: %w", err)
	}
	return page.Content, nil
}

// transferFunds transfers a specified amount of satoshis to the recipient's paymail address.
// It accepts a context parameter to manage cancellation and timeouts.
// On success, it returns the transaction representing the fund transfer and a nil error.
// If the API call fails, it returns a non-nil error.
func (l *leaderService) transferFunds(ctx context.Context, paymail string, funds uint64) (*response.Transaction, error) {
	return l.transferService.transferFunds(ctx, paymail, funds)
}

// createActor creates a new actor and registers its paymail within the leader service instance.
// It accepts a context parameter to manage cancellation and timeouts, as well as an alias to serve as a unique name for the actor.
// On success, it returns the registered actor instance and a nil error.
// If the API call fails, or if the actor's xPub or paymail registration fails, it returns a non-nil error.
func (l *leaderService) createActor(ctx context.Context, alias string) (*actor, error) {
	actor, err := newActor(alias, l.domain)
	if err != nil {
		return nil, fmt.Errorf("failed to create actor with alias %q: %w", alias, err)
	}

	_, err = l.adminAPI.CreateXPub(ctx, &commands.CreateUserXpub{XPub: actor.xPub})
	if err != nil {
		return nil, fmt.Errorf("failed to create xPub for actor %q: %w", alias, err)
	}

	_, err = l.adminAPI.CreatePaymail(ctx, &commands.CreatePaymail{
		Key:        actor.xPub,
		Address:    actor.paymail,
		PublicName: "Regression tests",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create paymail for actor %q: %w", alias, err)
	}

	l.actors = append(l.actors, actor)
	return actor, nil
}

// removeActors removes all actors and their associated paymails created by the leader service.
// It accepts a context parameter to manage cancellation and timeouts.
// On success, it returns a nil error, indicating that the removal process was completed successfully.
// If the API call fails, or if any paymail deletion fails, it returns an error group that wraps the errors encountered.
func (l *leaderService) removeActors(ctx context.Context) error {
	var errs []error
	for _, a := range l.actors {
		if err := l.adminAPI.DeletePaymail(ctx, a.paymail); err != nil {
			errs = append(errs, fmt.Errorf("failed to delete paymail %s: %w", a.paymail, err))
		}
	}
	l.actors = nil // Clear the list of actors.
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

// newLeaderService initializes and returns a new LeaderService instance.
// It accepts the configuration required for the service, including API URLs and xPriv keys.
// On success, it returns the initialized instance and a nil error.
// If there are issues with API initialization or domain retrieval, it returns a non-nil error.
func newLeaderService(cfg *leaderServiceConfig) (*leaderService, error) {
	walletCfg := config.New(config.WithAddr(cfg.envURL))
	adminAPI, err := wallet.NewAdminAPIWithXPriv(walletCfg, cfg.adminXPriv)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize admin API: %w", err)
	}
	userAPI, err := wallet.NewUserAPIWithXPriv(walletCfg, cfg.envXPriv)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize user API: %w", err)
	}
	transferService, err := newTransferService(userAPI)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize transfer service: %w", err)
	}

	sharedCfg, err := userAPI.SharedConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve paymail domain: %w", err)
	}
	if len(sharedCfg.PaymailDomains) != 1 {
		return nil, fmt.Errorf("expected one paymail domain, found %d", len(sharedCfg.PaymailDomains))
	}

	return &leaderService{
		cfg:             cfg,
		adminAPI:        adminAPI,
		userAPI:         userAPI,
		transferService: transferService,
		domain:          sharedCfg.PaymailDomains[0],
	}, nil
}
