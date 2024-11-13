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

func ExpectedCreatedAccessKey(t *testing.T) *response.AccessKey {
	return &response.AccessKey{
		Model: response.Model{
			Metadata: map[string]interface{}{
				"key": "value",
			},
			CreatedAt: ParseTime(t, "2024-11-13T11:44:04.95481Z"),
			UpdatedAt: ParseTime(t, "2024-11-13T12:44:04.954844+01:00"),
		},
		ID:     "d8558b86-9382-4c42-8ebe-7cca5d8de60b",
		XpubID: "345cef2e-36a7-4c28-b0a7-948bfdb03e5e",
		Key:    "dbb23e77-0467-4262-a0ef-3d30653866ae",
	}
}

func ExpectedRertrivedAccessKey(t *testing.T) *response.AccessKey {
	return &response.AccessKey{
		Model: response.Model{
			Metadata: map[string]interface{}{
				"key": "value",
			},
			CreatedAt: ParseTime(t, "2024-11-13T11:44:04.95481Z"),
			UpdatedAt: ParseTime(t, "2024-11-13T11:44:04.954844Z"),
		},
		ID:     "1fb70cc2-e9d9-41a3-842e-f71cc58d9787",
		XpubID: "e8d7d52f-01a1-4466-87fe-25a2225ef5e4",
	}
}

func ExpectedRevokedAccessKey(t *testing.T) *response.AccessKey {
	ts := ParseTime(t, "2024-11-13T12:54:36.987563+01:00")
	return &response.AccessKey{
		Model: response.Model{
			Metadata: map[string]interface{}{
				"key": "value",
			},
			CreatedAt: ParseTime(t, "2024-11-13T11:44:04.95481Z"),
			UpdatedAt: ParseTime(t, "2024-11-13T11:54:36.988715Z"),
		},
		ID:        "081743f7-040e-45a3-8c36-dde39001e20d",
		XpubID:    "41d0a43c-1721-4777-ad4a-57cbb2b38160",
		RevokedAt: &ts,
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
