package configs_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/configs"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/stretchr/testify/require"
)

func TestAPI_SharedConfig(t *testing.T) {
	type HTTPGetter interface {
		Get(ctx context.Context, path string) ([]byte, error)
	}
	tests := map[string]struct {
		HTTP             HTTPGetter
		expectedErr      error
		expectedResponse *response.SharedConfig
	}{
		"shared config: HTTP GET /api/v1/configs/shared response: 200": {
			HTTP: &SharedConfigAlwaysPass{
				ExpectedResponse: &response.SharedConfig{
					PaymailDomains: []string{
						"john.doe.test.4chain.space",
					},
					ExperimentalFeatures: map[string]bool{
						"pikeContactsEnabled": true,
						"pikePaymentEnabled":  false,
					},
				},
			},
			expectedResponse: &response.SharedConfig{
				PaymailDomains: []string{
					"john.doe.test.4chain.space",
				},
				ExperimentalFeatures: map[string]bool{
					"pikeContactsEnabled": true,
					"pikePaymentEnabled":  false,
				},
			},
		},
		"shared config: HTTP GET /api/v1/configs/shared response: 500": {
			HTTP: &SharedConfigAlwaysFailure{
				ExpectedErr: models.SPVError{
					Message:    http.StatusText(http.StatusInternalServerError),
					StatusCode: http.StatusInternalServerError,
					Code:       models.UnknownErrorCode,
				},
			},
			expectedErr: models.SPVError{
				Message:    http.StatusText(http.StatusInternalServerError),
				StatusCode: http.StatusInternalServerError,
				Code:       models.UnknownErrorCode,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			addr := "http://localhost:3003"
			API := configs.NewAPI(addr, tc.HTTP)
			ctx := context.Background()

			got, err := API.SharedConfig(ctx)
			require.ErrorIs(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedResponse, got)
		})
	}
}

// SharedConfigAlwaysPass stubs the HTTP client to simulate a successful response from
// the SPV Wallet API's `/api/v1/configs/shared` endpoint. This is useful for testing
// scenarios where a successful, predefined configuration response is needed.
type SharedConfigAlwaysPass struct {
	ExpectedResponse *response.SharedConfig
}

// Get simulates an HTTP GET request to the `/api/v1/configs/shared` endpoint and returns
// a successful response with no error. It marshals a predefined configuration into JSON format,
// mimicking the API's successful response.
func (s *SharedConfigAlwaysPass) Get(ctx context.Context, path string) ([]byte, error) {
	bb, err := json.Marshal(s.ExpectedResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal error stub response: %w", err)
	}
	return bb, nil
}

// SharedConfigAlwaysFailure stubs the HTTP client to simulate a failure response from
// the SPV Wallet API's `/api/v1/configs/shared` endpoint. This is useful for testing
// scenarios where a predefined error response is needed to verify error handling.
type SharedConfigAlwaysFailure struct {
	ExpectedErr error
}

// Get simulates an HTTP GET request to the `/api/v1/configs/shared` endpoint and returns
// a simulated failure response, mimicking the error structure expected from the SPV Wallet API.
func (s *SharedConfigAlwaysFailure) Get(ctx context.Context, path string) ([]byte, error) {
	return nil, s.ExpectedErr
}
