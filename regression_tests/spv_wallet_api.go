package regressiontests

import (
	"fmt"
)

// spvWalletConfig contains configuration settings for initializing a SPVWalletAPI instance.
// These include the environment URL and private keys required for admin and user operations.
type spvWalletConfig struct {
	envURL     string // URL of the SPV Wallet API environment.
	envXPriv   string // Extended private key (xPriv) for the user account.
	adminXPriv string // Extended private key (xPriv) for the admin account.
	name       string
}

// Validate validates the spvWalletConfig.
// It ensures that required fields like EnvURL and keys are not empty.
func (c *spvWalletConfig) Validate() error {
	if c.envURL == "" {
		return fmt.Errorf("validation failed: environment URL is required")
	}
	if c.adminXPriv == "" {
		return fmt.Errorf("validation failed: admin xPriv is required")
	}
	if c.envXPriv == "" {
		return fmt.Errorf("validation failed: leader xPriv is required")
	}
	if c.name == "" {
		return fmt.Errorf("validation failed: SPV Wallet API instance name is required")
	}
	return nil
}

// spvWallet represents the core API for interacting with the SPV Wallet ecosystem.
// It holds configuration and client instances for admin, user, and leader operations.
type spvWallet struct {
	cfg    *spvWalletConfig // Configuration for the SPV Wallet API (e.g., environment URL, keys).
	admin  *admin           // Admin client for performing administrative tasks like creating xPubs and paymails.
	user   *user            // User client for standard wallet operations, such as transactions and balance retrieval.
	leader *user            // Leader user client with potentially elevated privileges, managing broader wallet operations.
	name   string           // Name of the SPV Wallet API instance.
}

// initSPVWallet initializes the SPVWalletAPI with Admin, Leader, and User clients.
// It accepts an SPVWalletConfig and a user alias to be created as input parameters.
// On success, it returns an initialized SPVWalletAPI instance and nil error.
// If initialization of any component fails, a non-nil error is returned.
func initSPVWallet(cfg *spvWalletConfig, alias string) (*spvWallet, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	admin, err := initAdmin(cfg.envURL, cfg.adminXPriv)
	if err != nil {
		return nil, fmt.Errorf("could not initialize admin: %w", err)
	}

	leader, err := initLeaderUser(cfg.envURL, cfg.envXPriv)
	if err != nil {
		return nil, fmt.Errorf("could not initialize leader user: %w", err)
	}

	user, err := initUser(alias, cfg.envURL)
	if err != nil {
		return nil, fmt.Errorf("could not initialize user %q: %w", alias, err)
	}

	return &spvWallet{
		cfg:    cfg,
		admin:  admin,
		leader: leader,
		user:   user,
	}, nil
}
