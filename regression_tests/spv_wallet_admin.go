package regressiontests

import (
	"context"
	"errors"
	"fmt"

	wallet "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/config"
)

// admin represents an administrator within the SPV Wallet ecosystem.
// It includes the administrator's private key (xPriv) and provides access
// to the SPV Wallet's AdminAPI client for managing xPub and paymail-related operations.
type admin struct {
	xPriv    string           // The extended private key for the administrator.
	client   *wallet.AdminAPI // The API client for interacting with administrative functionalities in the SPV Wallet.
	paymails []string         // The paymail addresses created by administrator.
}

// createPaymail registers a new paymail address associated with the provided xPub key.
// It accepts a context parameter to manage cancellation and timeouts.
// The method uses the AdminAPI client to create the paymail record in the SPV Wallet ecosystem.
// On success, the paymail is added to the Admin's list of paymails, and nil is returned.
// If the operation fails, a non-nil error detailing the failure is returned.
func (a *admin) createPaymail(ctx context.Context, xPub, paymail string) error {
	_, err := a.client.CreatePaymail(ctx, &commands.CreatePaymail{
		Key:        xPub,
		Address:    paymail,
		PublicName: "Regression tests",
	})
	if err != nil {
		return fmt.Errorf("could not create paymail for actor %q: %w", paymail, err)
	}
	a.paymails = append(a.paymails, paymail)
	return nil
}

// createXPub registers a new xPub in the SPV Wallet ecosystem.
// It accepts a context parameter to manage cancellation and timeouts.
// The method uses the AdminAPI client to create the xPub record in the SPV Wallet ecosystem.
// On success, nil is returned. If the operation fails, a non-nil error detailing the failure is returned.
func (a *admin) createXPub(ctx context.Context, xPub string) error {
	_, err := a.client.CreateXPub(ctx, &commands.CreateUserXpub{XPub: xPub})
	if err != nil {
		return fmt.Errorf("could not create xPub record %s: %w", xPub, err)
	}
	return nil
}

// deletePaymails deletes all associated paymails created by the SPV Wallet API.
// It accepts a context parameter to manage cancellation and timeouts.
// On success, it returns a nil error, indicating that the removal process was completed successfully.
// If the API call fails, or if any paymail deletion fails, it returns an error group that wraps the errors encountered.
func (a *admin) deletePaymails(ctx context.Context) error {
	var errs []error
	for _, p := range a.paymails {
		if err := a.client.DeletePaymail(ctx, p); err != nil {
			errs = append(errs, fmt.Errorf("could not delete paymail %s: %w", p, err))
		}
	}
	a.paymails = nil // Clear the list of paymails.
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

// initAdmin initializes a new admin within the SPV Wallet ecosystem.
// It accepts the SPV Wallet API URL and the administrator's extended private key (xPriv) as input parameters.
// The function initializes the wallet's AdminAPI client using the provided xPriv,
// enabling the management of xPub and paymail-related operations.
// On success, it returns the initialized admin and a nil error.
// If the initialization fails, it returns a non-nil error with details of the failure.
func initAdmin(url, xPriv string) (*admin, error) {
	cfg := config.New(config.WithAddr(url))
	client, err := wallet.NewAdminAPIWithXPriv(cfg, xPriv)
	if err != nil {
		return nil, fmt.Errorf("could not initialize admin API: %w", err)
	}
	return &admin{xPriv: xPriv, client: client}, nil
}
