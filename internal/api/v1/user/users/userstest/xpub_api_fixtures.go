package userstest

import (
	"net/http"
	"testing"
	"time"

	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

func ExpectedUpdatedXPubMetadata(t *testing.T) *response.Xpub {
	return &response.Xpub{
		Model: response.Model{
			CreatedAt: ParseTime(t, "2024-10-07T13:39:07.886862Z"),
			UpdatedAt: ParseTime(t, "2024-11-13T11:41:56.115402Z"),
			Metadata: map[string]any{
				"metadata": map[string]any{
					"key": "value",
				},
			},
		},
		ID:              "1356cc11-8164-4364-afa4-59f096a79eb5",
		CurrentBalance:  315,
		NextInternalNum: 13,
		NextExternalNum: 2,
	}
}

func ExpectedUserXPub(t *testing.T) *response.Xpub {
	return &response.Xpub{
		Model: response.Model{
			CreatedAt: ParseTime(t, "2024-10-07T13:39:07.886862Z"),
			UpdatedAt: ParseTime(t, "2024-11-12T11:31:07.741621Z"),
			Metadata: map[string]any{
				"metadata": map[string]any{
					"key": "value",
				},
			},
		},
		ID:              "af64633f-b2ce-441e-9d61-acda0884eb53",
		CurrentBalance:  315,
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

func Ptr[T any](value T) *T {
	return &value
}