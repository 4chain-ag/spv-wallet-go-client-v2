package regressiontests

import (
	"context"
	"fmt"

	"github.com/bitcoin-sv/spv-wallet/models/response"
)

type Transactions []*response.Transaction

func (tt Transactions) Has(id string) bool {
	for _, t := range tt {
		if t.ID == id {
			return true
		}
	}
	return false
}

type TransferSender interface {
	Paymail() string
	TransferFunds(ctx context.Context, paymail string, funds uint64) (*response.Transaction, error)
	Balance(ctx context.Context) (uint64, error)
	Transactions(ctx context.Context) ([]*response.Transaction, error)
}

type TransferRecipient interface {
	Paymail() string
	Balance(ctx context.Context) (uint64, error)
	Transactions(ctx context.Context) ([]*response.Transaction, error)
}

type TransferService struct {
	Sender        TransferSender
	Recipient     TransferRecipient
	TransferFunds uint64
}

func (t *TransferService) Do(ctx context.Context) error {
	sender := t.Sender.Paymail()
	recipient := t.Recipient.Paymail()

	prevSenderBalance, err := t.Sender.Balance(ctx)
	if err != nil {
		return fmt.Errorf("Sender [paymail: %s] could not fetch balance before transaction: %w", sender, err)
	}
	prevRecipientBalance, err := t.Recipient.Balance(ctx)
	if err != nil {
		return fmt.Errorf("Recipient [paymail: %s] could not fetch balance before transaction: %w", recipient, err)
	}

	transaction, err := t.Sender.TransferFunds(ctx, recipient, t.TransferFunds)
	if err != nil {
		return fmt.Errorf("Could not transfer: %d funds from Sender [paymail: %s] to Recipient [paymail: %s]: %w", t.TransferFunds, sender, recipient, err)
	}
	currentSenderBalance, err := t.Sender.Balance(ctx)
	if err != nil {
		return fmt.Errorf("Sender [paymail: %s] could not fetch balance after transaction: %w", sender, err)
	}
	expectedSenderBalance := prevSenderBalance - t.TransferFunds - transaction.Fee
	if currentSenderBalance != expectedSenderBalance {
		return fmt.Errorf("Sender [paymail: %s] balance should be equal to: %d. Got: %d", sender, expectedSenderBalance, currentSenderBalance)
	}

	transactions, err := t.Recipient.Transactions(ctx)
	if err != nil {
		return fmt.Errorf("Recipient [paymail: %s] could not fetch transactions list: %w", recipient, err)
	}
	if !Transactions(transactions).Has(transaction.ID) {
		return fmt.Errorf("Recipient [paymail: %s] transactions list should contain transaction with ID: %s", recipient, transaction.ID)
	}
	currentRecipientBalance, err := t.Recipient.Balance(ctx)
	if err != nil {
		return fmt.Errorf("Recipient [paymail: %s] could not fetch balance after transaction: %w", recipient, err)
	}
	expectedRecipientBalance := prevRecipientBalance + t.TransferFunds
	if currentRecipientBalance != expectedRecipientBalance {
		return fmt.Errorf("Recipient [paymail: %s] balance should be equal to: %d. Got: %d", recipient, expectedRecipientBalance, currentSenderBalance)
	}
	return nil
}
