package userstest

import (
	"net/http"
	"testing"
	"time"

	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
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

func ExpectedAccessKeyPage(t *testing.T) *queries.AccessKeyPage {
	ts1 := ParseTime(t, "2024-11-13T11:54:36.987563Z")
	ts2 := ParseTime(t, "2024-11-08T13:43:18.599995Z")
	return &queries.AccessKeyPage{
		Content: []*response.AccessKey{
			{
				Model: response.Model{
					Metadata: map[string]interface{}{
						"key_1": "value_1",
					},
					CreatedAt: ParseTime(t, "2024-11-13T11:44:04.95481Z"),
					UpdatedAt: ParseTime(t, "2024-11-13T11:54:36.988715Z"),
				},
				ID:        "1f0504cd-d42d-4334-a441-a88a53aa47f8",
				XpubID:    "b271ae7e-ab17-4504-94c1-3a888f8b042a",
				RevokedAt: &ts1,
			},
			{
				Model: response.Model{
					Metadata: map[string]interface{}{
						"key_2": "value_2",
					},
					CreatedAt: ParseTime(t, "2024-11-13T11:07:43.595835Z"),
					UpdatedAt: ParseTime(t, "2024-11-13T11:07:43.595876Z"),
				},
				ID:     "41943e46-6999-409e-8dfd-d36ee75f1702",
				XpubID: "3e32dd04-72bd-4cc5-92da-123c29708472",
			},
			{
				Model: response.Model{
					Metadata: map[string]interface{}{
						"key_3": "value_3",
					},
					CreatedAt: ParseTime(t, "2024-11-08T13:43:18.554228Z"),
					UpdatedAt: ParseTime(t, "2024-11-08T13:43:18.60036Z"),
				},
				ID:        "41a87305-88f9-4d86-91f8-b2401078aaf9",
				XpubID:    "a035a7f0-2381-4d45-8a2d-197dd961f031",
				RevokedAt: &ts2,
			},
		},
		Page: response.PageDescription{
			Size:          50,
			Number:        1,
			TotalElements: 7,
			TotalPages:    1,
		},
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
