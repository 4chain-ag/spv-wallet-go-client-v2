package regressiontests

import (
	"fmt"

	wallet "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/config"
	"github.com/bitcoin-sv/spv-wallet-go-client/walletkeys"
)

// leaderConfig contains configuration settings for initializing a Leader.
// These include the environment URL and private keys required for admin and user operations.
type leaderConfig struct {
	envURL     string // URL of the SPV Wallet API environment.
	envXPriv   string // Extended private key (xPriv) for the user account.
	adminXPriv string // Extended private key (xPriv) for the admin account.
}

// leader represents an entity with administrative and user-level access to the SPV Wallet API.
// It manages wallet operations and maintains a list of associated paymails.
type leader struct {
	adminAPI *wallet.AdminAPI // adminAPI provides administrative access to perform privileged operations within the SPV Wallet.
	userAPI  *wallet.UserAPI  // userAPI provides user-level access for general wallet operations, such as fetching balances and transactions.
	domain   string           // domain represents the domain associated with the Leader's paymails.
	actor    *actor           // actor represents the details associated with the Leader's created actor.
}

// name returns the full name of the Leader, formatted as "Leader@<domain>".
func (l *leader) name() string { return "Leader@" + l.domain }

// setActor sets the actor associated with the Leader.
func (l *leader) setActor(a *actor) { l.actor = a }

// setPaymailDomain updates the domain associated with the Leader's paymails.
func (l *leader) setPaymailDomain(s string) { l.domain = s }

// NewLeader initializes and returns a new Leader instance for testing purposes.
// It sets up both the user and admin APIs using the provided configuration.
// The function fails the test immediately if any of the API initializations encounter an error.
func NewLeader(cfg *leaderConfig) (*leader, error) {
	walletCfg := config.New(config.WithAddr(cfg.envURL))
	userAPI, err := wallet.NewUserAPIWithXPriv(walletCfg, cfg.envXPriv)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize user API: %w", err)
	}
	adminAPI, err := wallet.NewAdminAPIWithXPriv(walletCfg, cfg.adminXPriv)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize admin API: %w", err)
	}
	return &leader{adminAPI: adminAPI, userAPI: userAPI}, nil
}

// actor represents an individual user or entity within the SPV Wallet ecosystem.
// It includes details like the alias, private key (xPriv), public key (xPub), and paymail address.
type actor struct {
	alias   string // The unique alias for the actor.
	xPriv   string // The extended private key for the actor.
	xPub    string // The extended public key for the actor.
	paymail string // The paymail address associated with the actor.
}

// createUserXpubCommand creates and returns a command to create a user xPub using the actor's xPub key.
func (a *actor) createUserXpubCommand() *commands.CreateUserXpub {
	return &commands.CreateUserXpub{XPub: a.xPub}
}

// createPaymailCommand creates and returns a command to create a paymail associated with the actor's xPub key.
// It uses the actor's paymail address and sets a public name for the paymail.
func (a *actor) createPaymailCommand() *commands.CreatePaymail {
	return &commands.CreatePaymail{
		Key:        a.xPub,
		Address:    a.paymail,
		PublicName: "Regression tests",
	}
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
