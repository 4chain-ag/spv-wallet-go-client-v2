package commands

// UpdateUserInformationMetadata contains the parameters needed to update a user's information metadata.
type UpdateUserInformationMetadata struct {
	Metadata map[string]any `json:"metadata"` // Metadata associated with the user.
}
