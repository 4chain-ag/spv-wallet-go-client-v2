package userstest

import (
	"net/http"
	"testing"
	"time"

	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

func ExpectedUserInformationUpdateMetadata(t *testing.T) *response.Xpub {
	return &response.Xpub{
		Model: response.Model{
			CreatedAt: ParseTime(t, "2024-10-07T13:39:07.886862Z"),
			UpdatedAt: ParseTime(t, "2024-11-11T15:46:10.050218Z"),
			Metadata: map[string]any{
				"metadata": map[string]any{
					"example_key": "example_value",
				},
			},
		},
		ID:              "98b5d2c8-c535-4cd1-9df1-7baa42474870",
		CurrentBalance:  215,
		NextInternalNum: 13,
		NextExternalNum: 2,
	}
}

func ExpectedUserInformation(t *testing.T) *response.Xpub {
	return &response.Xpub{
		Model: response.Model{
			CreatedAt: ParseTime(t, "2024-10-07T13:39:07.886862Z"),
			UpdatedAt: ParseTime(t, "2024-11-08T13:40:55.595522Z"),
			Metadata: map[string]any{
				"metadata": map[string]any{
					"some_metadata_2": "example2",
				},
			},
		},
		ID:              "dc7004d9-548b-42dd-a587-535fa456563f",
		CurrentBalance:  215,
		NextInternalNum: 13,
		NextExternalNum: 2,
	}
}

func NewBadRequestSPVError() *models.SPVError {
	return &models.SPVError{
		Message:    http.StatusText(http.StatusBadRequest),
		StatusCode: http.StatusBadRequest,
		Code:       "invalid-data-format",
	}
}

func ParseTime(t *testing.T, s string) time.Time {
	ts, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		t.Fatalf("test helper - time parse: %s", err)
	}
	return ts
}
