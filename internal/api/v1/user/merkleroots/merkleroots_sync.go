package merkleroots

import (
	"context"
	"errors"
	"fmt"

	"github.com/bitcoin-sv/spv-wallet/models"

	goclienterr "github.com/bitcoin-sv/spv-wallet-go-client/errors"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
)

// MerkleRootsRepository is an interface responsible for storing synchronized MerkleRoots and retrieving the last evaluation key from the database.
type MerkleRootsRepository interface {
	// GetLastMerkleRoot should return the Merkle root with the highest height from your memory, or undefined if empty.
	GetLastMerkleRoot() string
	// SaveMerkleRoots should store newly synced merkle roots into your storage;
	// NOTE: items are sorted in ascending order by block height.
	SaveMerkleRoots(syncedMerkleRoots []models.MerkleRoot) error
}

// SyncMerkleRoots syncs merkleroots known to spv-wallet with the client database
// If timeout is needed pass context.WithTimeout() as ctx param
// SyncMerkleRoots synchronizes Merkle roots known to SPV Wallet with the client database.
func (a *API) SyncMerkleRoots(ctx context.Context, repo MerkleRootsRepository) error {

	lastEvaluatedKey := repo.GetLastMerkleRoot()
	previousLastEvaluatedKey := lastEvaluatedKey

	for {
		select {
		case <-ctx.Done():
			return goclienterr.ErrSyncMerkleRootsTimeout
		default:
			// Query the MerkleRoots API
			result, err := a.MerkleRoots(ctx, queries.MerkleRootsQueryWithLastEvaluatedKey(lastEvaluatedKey))
			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					return goclienterr.ErrSyncMerkleRootsTimeout
				}
				return fmt.Errorf("failed to fetch merkle roots from API: %w", err)
			}

			// Handle empty results
			if len(result.Content) == 0 {
				return nil
			}

			// Update the last evaluated key
			lastEvaluatedKey = result.Page.LastEvaluatedKey
			if lastEvaluatedKey != "" && previousLastEvaluatedKey == lastEvaluatedKey {
				return goclienterr.ErrStaleLastEvaluatedKey
			}

			// Save fetched Merkle roots
			err = repo.SaveMerkleRoots(result.Content)
			if err != nil {
				return fmt.Errorf("failed to save merkle roots: %w", err)
			}

			if lastEvaluatedKey == "" {
				return nil
			}

			previousLastEvaluatedKey = lastEvaluatedKey
		}
	}
}
