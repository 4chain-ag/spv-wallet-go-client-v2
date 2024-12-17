package regressiontests

import (
	"context"
	"fmt"

	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

// fundsSender defines the methods required for a sender in a fund transfer operation.
// It includes functionality to retrieve the sender's xPub key for balance extraction
// and to send funds to specified recipients.
type fundsSender interface {
	XPub(ctx context.Context) (*response.Xpub, error)
	SendToRecipients(ctx context.Context, cmd *commands.SendToRecipients) (*response.Transaction, error)
}

// transferService provides an interface to interact with the SPV Wallet API.
// It supports operations such as retrieving balances and transferring funds within the wallet ecosystem.
// The service uses a FundsSender to perform actions related to fund transfers.
type transferService struct {
	sender fundsSender
}

// balance retrieves the current satoshi balance for the sender.
// It accepts a context parameter to manage cancellation and timeouts.
// On success, it returns the current balance and a nil error.
// If the API call fails, it returns a non-nil error with details of the failure.
func (t *transferService) balance(ctx context.Context) (uint64, error) {
	xpub, err := t.sender.XPub(ctx)
	if err != nil {
		return 0, fmt.Errorf("could not fetch xPub to retrieve current balance: %w", err)
	}
	return xpub.CurrentBalance, nil
}

// transferFunds sends a specified amount of satoshis to a recipient's paymail.
// It accepts a context parameter to manage cancellation and timeouts.
// On success, it returns the transaction representing the fund transfer and a nil error.
// If the actor has insufficient funds or the API call fails, it returns a non-nil error.
func (t *transferService) transferFunds(ctx context.Context, paymail string, funds uint64) (*response.Transaction, error) {
	balance, err := t.balance(ctx)
	if err != nil {
		return nil, fmt.Errorf("balance failed: %w", err)
	}
	if balance < funds {
		return nil, fmt.Errorf("insufficient balance: %d available, %d required", balance, funds)
	}

	recipient := commands.Recipients{To: paymail, Satoshis: funds}
	transaction, err := t.sender.SendToRecipients(ctx, &commands.SendToRecipients{
		Recipients: []*commands.Recipients{&recipient},
		Metadata:   map[string]any{"description": "regression-test"},
	})
	if err != nil {
		return nil, fmt.Errorf("could not transfer funds to %s: %w", paymail, err)
	}
	return transaction, nil
}

// newTransferService initializes and returns a new TransferService instance.
// It accepts a FundsSender, which provides the necessary functionality for retrieving wallet information
// and transferring funds. If the provided FundsSender is nil, an error is returned.
func newTransferService(f fundsSender) (*transferService, error) {
	if f == nil {
		return nil, fmt.Errorf("could not initialize transfer service: nil funds sender specified")
	}
	return &transferService{sender: f}, nil
}
