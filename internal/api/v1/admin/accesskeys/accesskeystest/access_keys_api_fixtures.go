package accesskeystest

import (
	"net/http"
	"testing"
	"time"

	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

func ExpectedAccessKeyPage(t *testing.T) *queries.AccessKeyPage {
	return &queries.AccessKeyPage{
		Content: []*response.AccessKey{
			{
				Model: response.Model{
					CreatedAt: parseTime(t, "2024-11-28T14:56:59.841638Z"),
					UpdatedAt: parseTime(t, "2024-11-28T14:56:59.841832Z"),
				},
				ID:     "3a77c921-b881-4907-8dc6-3903700272cb",
				XpubID: "cd6709cd-4f0e-464b-8d7d-0197e853f375",
			},
			{
				Model: response.Model{
					CreatedAt: parseTime(t, "2024-11-28T13:28:22.943632Z"),
					UpdatedAt: parseTime(t, "2024-11-28T13:28:22.943664Z"),
				},
				ID:     "35aacdfd-5839-4125-9180-d33e798b1cde",
				XpubID: "7c6c4462-626c-47f6-84bc-04044798a4bf",
			},
		},
		Page: response.PageDescription{
			Size:          2,
			Number:        1,
			TotalElements: 2,
			TotalPages:    1,
		},
	}
}

func NewBadRequestSPVError() models.SPVError {
	return models.SPVError{
		Message:    http.StatusText(http.StatusBadRequest),
		StatusCode: http.StatusBadRequest,
		Code:       "invalid-data-format",
	}
}

func NewInternalServerSPVError() models.SPVError {
	return models.SPVError{
		Message:    http.StatusText(http.StatusInternalServerError),
		StatusCode: http.StatusInternalServerError,
		Code:       models.UnknownErrorCode,
	}
}

func parseTime(t *testing.T, s string) time.Time {
	ts, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		t.Fatalf("test helper - time parse: %s", err)
	}
	return ts
}
