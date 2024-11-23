package merkleroots

import "github.com/bitcoin-sv/spv-wallet/models"

type NoopDB struct{}

func (NoopDB) GetLastMerkleRoot() string { return "" }

func (NoopDB) SaveMerkleRoots([]models.MerkleRoot) error { return nil }

func (NoopDB) IsNoop() bool { return true }

type DBNoopExtended struct {
	db DB
}

func (d DBNoopExtended) GetLastMerkleRoot() string { return d.db.GetLastMerkleRoot() }

func (d DBNoopExtended) SaveMerkleRoots(mm []models.MerkleRoot) error {
	return d.db.SaveMerkleRoots(mm)
}

func (d DBNoopExtended) IsNoop() bool { return false }
