package commands

import (
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/querybuilders"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

// RecordTransaction holds the arguments required to record a user transaction.
type RecordTransaction struct {
	Metadata    querybuilders.Metadata `json:"metadata"`    // Metadata associated with the transaction.
	Hex         string                 `json:"hex"`         // Hexadecimal string representation of the transaction.
	ReferenceID string                 `json:"referenceId"` // Reference ID for the transaction.
}

// DraftTransaction holds the arguments required to create user draft transaction.
type DraftTransaction struct {
	Config   response.TransactionConfig `json:"config"`   // Configuration for the transaction.
	Metadata querybuilders.Metadata     `json:"metadata"` // Metadata related to the transaction.
}

// UpdateTransactionMetadata holds the arguments required to update a user transaction's metadata.
type UpdateTransactionMetadata struct {
	ID       string                 `json:"-"`        // Unique identifier of the transaction to be updated.
	Metadata querybuilders.Metadata `json:"metadata"` // New metadata to associate with the transaction.
}