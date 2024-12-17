package regressiontests

import (
	"context"
	"fmt"

	"github.com/bitcoin-sv/spv-wallet/models/response"
)

// transactionsSlice represents a slice of response.Transaction objects.
type transactionsSlice []*response.Transaction

// has checks if a transaction with the specified ID exists in the transactions slice.
// It returns true if a transaction with the given ID is found, and false otherwise.
func (tt transactionsSlice) has(id string) bool {
	for _, t := range tt {
		if t.ID == id {
			return true
		}
	}
	return false
}

// transferSender is an interface that defines methods for a sender in a fund transfer operation.
// It includes methods for retrieving the sender's paymail, balance, and transactions,
// as well as transferring funds to a recipient.
type transferSender interface {
	paymail() string
	transferFunds(ctx context.Context, paymail string, funds uint64) (*response.Transaction, error)
	balance(ctx context.Context) (uint64, error)
	transactions(ctx context.Context) (transactionsSlice, error)
}

// transferRecipient is an interface that defines methods for a recipient in a fund transfer operation.
// It includes methods for retrieving the recipient's paymail, balance, and transactions.
type transferRecipient interface {
	paymail() string
	balance(ctx context.Context) (uint64, error)
	transactions(ctx context.Context) (transactionsSlice, error)
}

// transferVerficiationService provides functionality to verify a transfer between a sender and recipient.
// It ensures that the sender's balance is reduced correctly and the recipient's balance increases appropriately,
// as well as verifies that the transaction appears in the recipient's transaction history.
type transferVerficiationService struct {
	sender        transferSender    // The sender involved in the transfer.
	recipient     transferRecipient // The recipient involved in the transfer.
	transferFunds uint64            // The amount of funds to be transferred.
}

// funds returns the amount of funds to be transferred as part of the verification process.
// It retrieves the value of the TransferFunds field, which represents the amount being transferred.
func (t *transferVerficiationService) funds() uint64 {
	return t.transferFunds
}

// senderPaymail returns the paymail address of the sender involved in the transfer verification.
// It calls the Paymail method of the Sender interface to retrieve the sender's paymail.
func (t *transferVerficiationService) senderPaymail() string {
	return t.sender.paymail()
}

// recipientPaymail returns the paymail address of the recipient involved in the transfer verification.
// It calls the Paymail method of the Recipient interface to retrieve the recipient's paymail.
func (t *transferVerficiationService) recipientPaymail() string {
	return t.recipient.paymail()
}

// do performs the transfer verification by checking balances and transaction history before and after the transfer.
// It verifies that the sender's balance is reduced by the transfer amount and transaction fee,
// and that the recipient's balance is increased by the transfer amount.
// Additionally, it ensures the transaction appears in the recipient's transaction list.
// If any step fails, it returns a non-nil error with a description of the issue.
func (t *transferVerficiationService) do(ctx context.Context) error {
	sender := t.senderPaymail()
	recipient := t.recipientPaymail()

	// Fetch previous balances for sender and recipient
	prevSenderBalance, err := t.sender.balance(ctx)
	if err != nil {
		return fmt.Errorf("Sender [paymail: %s] could not fetch balance before transaction: %w", sender, err)
	}
	prevRecipientBalance, err := t.recipient.balance(ctx)
	if err != nil {
		return fmt.Errorf("Recipient [paymail: %s] could not fetch balance before transaction: %w", recipient, err)
	}

	// Perform the transfer
	transaction, err := t.sender.transferFunds(ctx, recipient, t.transferFunds)
	if err != nil {
		return fmt.Errorf("Could not transfer: %d funds from Sender [paymail: %s] to Recipient [paymail: %s]: %w", t.transferFunds, sender, recipient, err)
	}

	// Verify sender's balance after the transfer
	currentSenderBalance, err := t.sender.balance(ctx)
	if err != nil {
		return fmt.Errorf("Sender [paymail: %s] could not fetch balance after transaction: %w", sender, err)
	}
	expectedSenderBalance := prevSenderBalance - t.transferFunds - transaction.Fee
	if currentSenderBalance != expectedSenderBalance {
		return fmt.Errorf("Sender [paymail: %s] balance should be equal to: %d. Got: %d", sender, expectedSenderBalance, currentSenderBalance)
	}

	// Verify recipient's transaction list
	transactions, err := t.recipient.transactions(ctx)
	if err != nil {
		return fmt.Errorf("Recipient [paymail: %s] could not fetch transactions list: %w", recipient, err)
	}
	if !transactions.has(transaction.ID) {
		return fmt.Errorf("Recipient [paymail: %s] transactions list should contain transaction with ID: %s", recipient, transaction.ID)
	}

	// Verify recipient's balance after the transfer
	currentRecipientBalance, err := t.recipient.balance(ctx)
	if err != nil {
		return fmt.Errorf("Recipient [paymail: %s] could not fetch balance after transaction: %w", recipient, err)
	}
	expectedRecipientBalance := prevRecipientBalance + t.transferFunds
	if currentRecipientBalance != expectedRecipientBalance {
		return fmt.Errorf("Recipient [paymail: %s] balance should be equal to: %d. Got: %d", recipient, expectedRecipientBalance, currentSenderBalance)
	}
	return nil
}
