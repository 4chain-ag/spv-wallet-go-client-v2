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

// LeaderServiceConfig contains configuration settings for initializing a LeaderService.
// These include the environment URL and private keys required for admin and user operations.
type LeaderServiceConfig struct {
	EnvURL     string // URL of the SPV Wallet API environment.
	EnvXPriv   string // Extended private key (xPriv) for the user account.
	AdminXPriv string // Extended private key (xPriv) for the admin account.
}

// LeaderService provides an interface to interact with the SPV Wallet API.
// It supports operations like retrieving balances, managing transactions,
// creating and managing actors, and transferring funds within the wallet ecosystem.
type LeaderService struct {
	cfg      *LeaderServiceConfig // Configuration for the service.
	adminAPI *wallet.AdminAPI     // Admin API client for privileged operations.
	userAPI  *wallet.UserAPI      // User API client for user-level operations.
	actors   []*Actor             // List of actors created by this service.
	domain   string               // Paymail domain used by the service.
}

// Paymail returns the paymail address of the leader.
// The paymail is constructed as "Leader@" followed by the configured domain.
func (l *LeaderService) Paymail() string {
	return "Leader@" + l.domain
}

// Balance retrieves the current satoshi balance for the leader account.
// It accepts a context parameter to manage cancellation and timeouts.
// The method fetches the xPub key of the user account and uses it to determine the current balance.
// On success, it returns the balance in satoshis and a nil error. If the API call fails, it returns a non-nil error.
func (l *LeaderService) Balance(ctx context.Context) (uint64, error) {
	xpub, err := l.userAPI.XPub(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch xPub for balance retrieval: %w", err)
	}
	return xpub.CurrentBalance, nil
}

// Transactions retrieves the list of transactions associated with the leader's account.
// It accepts a context parameter to manage cancellation and timeouts.
// On success, it returns a slice of transactions contained in the paginated response and a nil error.
// If the API call fails, it returns a non-nil error.
func (l *LeaderService) Transactions(ctx context.Context) ([]*response.Transaction, error) {
	page, err := l.userAPI.Transactions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transactions: %w", err)
	}
	return page.Content, nil
}

// TransferFunds transfers a specified amount of satoshis to the recipient's paymail address.
// It accepts a context parameter to manage cancellation and timeouts.
// On success, it returns the transaction representing the fund transfer and a nil error.
// If the API call fails, it returns a non-nil error.
func (l *LeaderService) TransferFunds(ctx context.Context, paymail string, funds uint64) (*response.Transaction, error) {
	balance, err := l.Balance(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch balance: %w", err)
	}
	if balance < funds {
		return nil, fmt.Errorf("insufficient balance: %d available, %d required", balance, funds)
	}

	recipient := commands.Recipients{To: paymail, Satoshis: funds}
	transaction, err := l.userAPI.SendToRecipients(ctx, &commands.SendToRecipients{
		Recipients: []*commands.Recipients{&recipient},
		Metadata:   map[string]any{"description": "regression-test"},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to transfer funds to %s: %w", paymail, err)
	}
	return transaction, nil
}

// CreateActor creates a new actor and registers its paymail within the leader service instance.
// It accepts a context parameter to manage cancellation and timeouts, as well as an alias to serve as a unique name for the actor.
// On success, it returns the registered actor instance and a nil error.
// If the API call fails, or if the actor's xPub or paymail registration fails, it returns a non-nil error.
func (l *LeaderService) CreateActor(ctx context.Context, alias string) (*Actor, error) {
	actor, err := NewActor(alias, l.domain)
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

// RemoveActors removes all actors and their associated paymails created by the leader service.
// It accepts a context parameter to manage cancellation and timeouts.
// On success, it returns a nil error, indicating that the removal process was completed successfully.
// If the API call fails, or if any paymail deletion fails, it returns an error group that wraps the errors encountered.
func (l *LeaderService) RemoveActors(ctx context.Context) error {
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

// NewLeaderService initializes and returns a new LeaderService instance.
// It accepts the configuration required for the service, including API URLs and xPriv keys.
// On success, it returns the initialized instance and a nil error.
// If there are issues with API initialization or domain retrieval, it returns a non-nil error.
func NewLeaderService(cfg *LeaderServiceConfig) (*LeaderService, error) {
	walletCfg := config.New(config.WithAddr(cfg.EnvURL))
	adminAPI, err := wallet.NewAdminAPIWithXPriv(walletCfg, cfg.AdminXPriv)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize admin API: %w", err)
	}

	userAPI, err := wallet.NewUserAPIWithXPriv(walletCfg, cfg.EnvXPriv)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize user API: %w", err)
	}

	sharedCfg, err := userAPI.SharedConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve paymail domain: %w", err)
	}
	if len(sharedCfg.PaymailDomains) != 1 {
		return nil, fmt.Errorf("expected one paymail domain, found %d", len(sharedCfg.PaymailDomains))
	}

	return &LeaderService{
		cfg:      cfg,
		adminAPI: adminAPI,
		userAPI:  userAPI,
		domain:   sharedCfg.PaymailDomains[0],
	}, nil
}
