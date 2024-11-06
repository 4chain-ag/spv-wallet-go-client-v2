package client

import (
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/querybuilders"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/transactions"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

// RecordTransactionArgs holds the arguments required to record a user transaction.
type RecordTransactionArgs struct {
	Metadata    querybuilders.Metadata // Metadata associated with the transaction.
	Hex         string                 // Hexadecimal string representation of the transaction.
	ReferenceID string                 // Reference ID for the transaction.
}

// ParseToRecordTransactionRequest converts RecordTransactionArgs to a RecordTransactionRequest
// for SPV Wallet API consumption.
func (r RecordTransactionArgs) ParseToRecordTransactionRequest() *transactions.RecordTransactionRequest {
	return &transactions.RecordTransactionRequest{
		Metadata:    r.Metadata,
		Hex:         r.Hex,
		ReferenceID: r.ReferenceID,
	}
}

// DraftTransactionArgs holds the arguments required to create user draft transaction.
type DraftTransactionArgs struct {
	Config   response.TransactionConfig // Configuration for the transaction.
	Metadata querybuilders.Metadata     // Metadata related to the transaction.
}

// ParseToDraftTransactionRequest converts DraftTransactionArgs to a DraftTransactionRequest
// for SPV Wallet API consumption.
func (d DraftTransactionArgs) ParseToDraftTransactionRequest() *transactions.DraftTransactionRequest {
	return &transactions.DraftTransactionRequest{
		Config:   d.Config,
		Metadata: d.Metadata,
	}
}

// UpdateTransactionMetadataArgs holds the arguments required to update a user transaction's metadata.
type UpdateTransactionMetadataArgs struct {
	ID       string                 // Unique identifier of the transaction to be updated.
	Metadata querybuilders.Metadata // New metadata to associate with the transaction.
}

// ParseUpdateTransactionMetadataRequest converts UpdateTransactionMetadataArgs to an
// UpdateTransactionMetadataRequest for SPV Wallet API consumption.
func (u UpdateTransactionMetadataArgs) ParseUpdateTransactionMetadataRequest() *transactions.UpdateTransactionMetadataRequest {
	return &transactions.UpdateTransactionMetadataRequest{
		ID:       u.ID,
		Metadata: u.Metadata,
	}
}
