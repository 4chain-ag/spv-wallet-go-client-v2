package merkleroots

import (
	"net/url"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/querybuilders"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
)

type merkleRootsFilterBuilder struct {
	query queries.MerkleRootsQuery
}

func (m *merkleRootsFilterBuilder) Build() (url.Values, error) {
	params := querybuilders.NewExtendedURLValues()
	params.AddPair("batchSize", m.query.BatchSize)
	params.AddPair("lastEvaluatedKey", m.query.LastEvaluatedKey)
	return params.Values, nil
}
