package regressiontests

import (
	"context"
	"fmt"
	"os"
	"testing"

	wallet "github.com/bitcoin-sv/spv-wallet-go-client"
)

// transferBalance represents a structure to verify the correctness of a balance transfer
// between a sender and a recipient in the SPV Wallet ecosystem.
type transferBalance struct {
	// Transaction participants
	sender    *user // The sender of the funds.
	recipient *user // The recipient of the funds.

	// Previous balances
	senderBalance    uint64 // The sender's balance before the transaction.
	recipientBalance uint64 // The recipient's balance before the transaction.

	// Transaction details
	transactionID string // The unique ID of the transaction.
	fee           uint64 // The transaction fee deducted from the sender.
	funds         uint64 // The amount of funds transferred.
}

// check validates that the transfer operation has been successfully executed.
// It checks the following:
// - The sender's balance has been correctly reduced by the transferred amount and fee.
// - The recipient's balance has been correctly increased by the transferred amount.
// - The transaction appears in the recipient's transaction list.
//
// If any validation fails, the test fails with a detailed error message.
func (tr *transferBalance) check(ctx context.Context, t *testing.T) {
	t.Helper()

	// Verify sender's balance after the transaction
	actualBalance := checkBalance(ctx, t, tr.sender.client)
	expectedBalance := tr.senderBalance - tr.funds - tr.fee
	if actualBalance != expectedBalance {
		t.Errorf("Transfer funds %d wasn't successful from sender %s to recipient %s. Expected sender balance to decrease from %d to %d, but got %d.",
			tr.funds, tr.sender.paymail, tr.recipient.paymail, tr.senderBalance, expectedBalance, actualBalance)
	}

	// Verify that the transaction appears in the recipient's transaction list
	page, err := tr.recipient.client.Transactions(ctx)
	if err != nil {
		t.Errorf("Failed to retrieve transactions for recipient %s. Expected nil error, got: %v", tr.recipient.paymail, err)
	}

	recipientTransactions := transactionsSlice(page.Content)
	if !recipientTransactions.Has(tr.transactionID) {
		t.Errorf("Transaction %s was not found in recipient %s's transaction list. Sent by %s.",
			tr.transactionID, tr.recipient.paymail, tr.sender.paymail)
	}

	// Verify recipient's balance after the transaction
	actualBalance = checkBalance(ctx, t, tr.recipient.client)
	expectedBalance = tr.recipientBalance + tr.funds
	if actualBalance != expectedBalance {
		t.Errorf("Transfer funds %d wasn't successful from sender %s to recipient %s. Expected recipient balance to increase from %d to %d, but got %d.",
			tr.funds, tr.sender.paymail, tr.recipient.paymail, tr.recipientBalance, expectedBalance, actualBalance)
	}
}

// checkBalance retrieves and returns the current balance for the given user client.
// If the balance retrieval fails, the test fails with a fatal error.
func checkBalance(ctx context.Context, t *testing.T, client *wallet.UserAPI) uint64 {
	t.Helper()

	// Fetch the extended public key (xPub) and associated balance
	xPub, err := client.XPub(ctx)
	if err != nil {
		t.Fatalf("Balance check failed: could not fetch the balance: %v", err)
	}
	if xPub == nil {
		t.Fatal("Balance check failed: expected non-nil xPub response.")
	}

	return xPub.CurrentBalance
}

// lookupEnvOrDefault retrieves the value of the specified environment variable.
// If the variable is not set, it returns the provided default value.
func lookupEnvOrDefault(t *testing.T, env string, defaultValue string) string {
	t.Helper()

	v, ok := os.LookupEnv(env)
	if !ok {
		fmt.Printf("Environment variable %s not set, using default: %s\n", env, defaultValue)
		return defaultValue
	}
	return v
}
