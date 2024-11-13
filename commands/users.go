package commands

// UpdateXPubMetadata contains the parameters needed to update a user's xpub metadata.
type UpdateXPubMetadata struct {
	Metadata map[string]any `json:"metadata"` // Metadata associated with the current user's xpub
}

// GenerateAccessKey contains the parameters needed to update a user's access key metadata.
type GenerateAccessKey struct {
	Metadata map[string]any `json:"metadata"` // Metadata associated with user's access key.
}
