package regressiontests

import (
	"context"
	"fmt"

	"github.com/bitcoin-sv/spv-wallet/models/response"
)

// Transactions represents a slice of response.Transaction objects.
type Transactions []*response.Transaction

// Has checks if a transaction with the specified ID exists in the Transactions slice.
// It returns true if a transaction with the given ID is found, and false otherwise.
func (tt Transactions) Has(id string) bool {
	for _, t := range tt {
		if t.ID == id {
			return true
		}
	}
	return false
}

// TransferSender is an interface that defines methods for a sender in a fund transfer operation.
// It includes methods for retrieving the sender's paymail, balance, and transactions,
// as well as transferring funds to a recipient.
type TransferSender interface {
	Paymail() string
	TransferFunds(ctx context.Context, paymail string, funds uint64) (*response.Transaction, error)
	Balance(ctx context.Context) (uint64, error)
	Transactions(ctx context.Context) ([]*response.Transaction, error)
}

// TransferRecipient is an interface that defines methods for a recipient in a fund transfer operation.
// It includes methods for retrieving the recipient's paymail, balance, and transactions.
type TransferRecipient interface {
	Paymail() string
	Balance(ctx context.Context) (uint64, error)
	Transactions(ctx context.Context) ([]*response.Transaction, error)
}

// TransferVerficiationService provides functionality to verify a transfer between a sender and recipient.
// It ensures that the sender's balance is reduced correctly and the recipient's balance increases appropriately,
// as well as verifies that the transaction appears in the recipient's transaction history.
type TransferVerficiationService struct {
	Sender        TransferSender    // The sender involved in the transfer.
	Recipient     TransferRecipient // The recipient involved in the transfer.
	TransferFunds uint64            // The amount of funds to be transferred.
}

// Funds returns the amount of funds to be transferred as part of the verification process.
// It retrieves the value of the TransferFunds field, which represents the amount being transferred.
func (t *TransferVerficiationService) Funds() uint64 {
	return t.TransferFunds
}

// SenderPaymail returns the paymail address of the sender involved in the transfer verification.
// It calls the Paymail method of the Sender interface to retrieve the sender's paymail.
func (t *TransferVerficiationService) SenderPaymail() string {
	return t.Sender.Paymail()
}

// RecipientPaymail returns the paymail address of the recipient involved in the transfer verification.
// It calls the Paymail method of the Recipient interface to retrieve the recipient's paymail.
func (t *TransferVerficiationService) RecipientPaymail() string {
	return t.Recipient.Paymail()
}

// Do performs the transfer verification by checking balances and transaction history before and after the transfer.
// It verifies that the sender's balance is reduced by the transfer amount and transaction fee,
// and that the recipient's balance is increased by the transfer amount.
// Additionally, it ensures the transaction appears in the recipient's transaction list.
// If any step fails, it returns a non-nil error with a description of the issue.
func (t *TransferVerficiationService) Do(ctx context.Context) error {
	sender := t.SenderPaymail()
	recipient := t.RecipientPaymail()

	// Fetch previous balances for sender and recipient
	prevSenderBalance, err := t.Sender.Balance(ctx)
	if err != nil {
		return fmt.Errorf("Sender [paymail: %s] could not fetch balance before transaction: %w", sender, err)
	}
	prevRecipientBalance, err := t.Recipient.Balance(ctx)
	if err != nil {
		return fmt.Errorf("Recipient [paymail: %s] could not fetch balance before transaction: %w", recipient, err)
	}

	// Perform the transfer
	transaction, err := t.Sender.TransferFunds(ctx, recipient, t.TransferFunds)
	if err != nil {
		return fmt.Errorf("Could not transfer: %d funds from Sender [paymail: %s] to Recipient [paymail: %s]: %w", t.TransferFunds, sender, recipient, err)
	}

	// Verify sender's balance after the transfer
	currentSenderBalance, err := t.Sender.Balance(ctx)
	if err != nil {
		return fmt.Errorf("Sender [paymail: %s] could not fetch balance after transaction: %w", sender, err)
	}
	expectedSenderBalance := prevSenderBalance - t.TransferFunds - transaction.Fee
	if currentSenderBalance != expectedSenderBalance {
		return fmt.Errorf("Sender [paymail: %s] balance should be equal to: %d. Got: %d", sender, expectedSenderBalance, currentSenderBalance)
	}

	// Verify recipient's transaction list
	transactions, err := t.Recipient.Transactions(ctx)
	if err != nil {
		return fmt.Errorf("Recipient [paymail: %s] could not fetch transactions list: %w", recipient, err)
	}
	if !Transactions(transactions).Has(transaction.ID) {
		return fmt.Errorf("Recipient [paymail: %s] transactions list should contain transaction with ID: %s", recipient, transaction.ID)
	}

	// Verify recipient's balance after the transfer
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
