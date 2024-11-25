package merkleroots

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	bip32 "github.com/bitcoin-sv/go-sdk/compat/bip32"
	goclienterr "github.com/bitcoin-sv/spv-wallet-go-client/errors"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/go-resty/resty/v2"
)

// MerkleRootsRepository is an interface responsible for storing synchronized MerkleRoots and retrieving the last evaluation key from the database.
type MerkleRootsRepository interface {
	// GetLastMerkleRoot should return the Merkle root with the highest height from your memory, or undefined if empty.
	GetLastMerkleRoot() string
	// SaveMerkleRoots should store newly synced merkle roots into your storage;
	// NOTE: items are sorted in ascending order by block height.
	SaveMerkleRoots(syncedMerkleRoots []models.MerkleRoot) error
}

type Client struct {
	merkleRootsAPI *API
	xPriv          *bip32.ExtendedKey
}

// NewClient initializes a new Merkle Roots Client.
// It requires the base URL of the API and a Resty HTTP client.
// The xPriv is optional and can be nil if authentication is not needed.
func NewClient(baseURL string, httpClient *resty.Client, xPriv *bip32.ExtendedKey) (*Client, error) {
	if baseURL == "" {
		return nil, errors.New("baseURL cannot be empty")
	}

	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	// Initialize the Merkle Roots API
	merkleRootsAPI := NewAPI(parsedURL, httpClient)

	return &Client{
		merkleRootsAPI: merkleRootsAPI,
		xPriv:          xPriv,
	}, nil
}

// SyncMerkleRoots syncs merkleroots known to spv-wallet with the client database
// If timeout is needed pass context.WithTimeout() as ctx param
// SyncMerkleRoots synchronizes Merkle roots known to SPV Wallet with the client database.
func (wc *Client) SyncMerkleRoots(ctx context.Context, repo MerkleRootsRepository) error {
	// Check if merkleRootsAPI is initialized
	if wc.merkleRootsAPI == nil {
		return errors.New("merkleRootsAPI is not initialized")
	}

	lastEvaluatedKey := repo.GetLastMerkleRoot()
	previousLastEvaluatedKey := lastEvaluatedKey

	for {
		select {
		case <-ctx.Done():
			return goclienterr.ErrSyncMerkleRootsTimeout
		default:
			// Query the MerkleRoots API
			result, err := wc.merkleRootsAPI.MerkleRoots(ctx, queries.MerkleRootsQueryWithLastEvaluatedKey(lastEvaluatedKey))
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
