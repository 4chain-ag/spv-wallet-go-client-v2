package restyutil_test

import (
	"errors"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/config"
	goclienterr "github.com/bitcoin-sv/spv-wallet-go-client/errors"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/restyutil"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

/*
 What is the test doing?
  - The test verifies the OnAfterResponse callback of NewHTTPClient.
  - The test verifies the behavior of the middleware when the response status code is 400.
  - The test verifies the behavior of the middleware when the response status code is 500.
  - The test verifies the behavior of the middleware when the response status code is 200.

  * Client Error (400)

  What happens in OnAfterResponse?
   - The middleware checks if the status code is 400.
   - It attempts to unmarshal the response body into models.SPVError.
   - Since the response body matches the expected format, it successfully unmarshals into a models.SPVError instance.
   - The middleware returns the unmarshaled SPVError as the error.

  * Server Error (500)

  What happens in OnAfterResponse?
   - The middleware detects a failure response (status code ≥ 400).
   - It attempts to unmarshal the response body into models.SPVError.
   - The unmarshaling fails (because the response is plain text, not JSON), and the middleware falls back to a generic error:

  * Success Response (200)

  What happens in OnAfterResponse?
   - The middleware detects a success response (IsSuccess() is true).
   - It returns nil, allowing the client to handle the response as usual.

 How Does the Test Verify OnAfterResponse?
 - Middleware Application: The OnAfterResponse middleware is invoked automatically because it is configured as part of the NewHTTPClient setup.

 Response Validation:
 - For successful responses, the test ensures no error is returned, confirming the middleware allowed the response to pass.
 - For error responses, the test ensures the correct error is returned, confirming the middleware processed the response as intended.
 - For unexpected formats, the test ensures fallback logic works, confirming the middleware gracefully handles unrecognized responses.
*/

// mockAuthenticator is a mock implementation of Authenticator interface
type mockAuthenticator struct{}

// Authenticate is a mock implementation of Authenticator interface
func (m *mockAuthenticator) Authenticate(r *resty.Request) error {
	return nil
}

// TestNewHTTPClient_OnAfterResponse tests the OnAfterResponse callback of NewHTTPClient
func TestNewHTTPClient_OnAfterResponse(t *testing.T) {
	tests := map[string]struct {
		statusCode       int
		responseBody     interface{}
		expectedError    error
		expectedSPVError *models.SPVError
	}{
		"Success Response 200": {
			statusCode:    200,
			responseBody:  map[string]string{"message": "success"},
			expectedError: nil,
		},
		"Client Error 400": {
			statusCode: 400,
			responseBody: models.SPVError{
				Message:    "Invalid request",
				StatusCode: 400,
				Code:       "invalid-request",
			},
			expectedError: models.SPVError{
				Message:    "Invalid request",
				StatusCode: 400,
				Code:       "invalid-request",
			},
		},
		"Server Error 500": {
			statusCode:   500,
			responseBody: "Internal server error",
			expectedError: models.SPVError{
				Message:    goclienterr.ErrUnrecognizedAPIResponse.Error(),
				StatusCode: 500,
				Code:       "internal-server-error",
			},
		},
	}

	// Mock configuration
	cfg := config.Config{
		Addr:      "http://mock-api",
		Timeout:   5,
		Transport: httpmock.DefaultTransport, // Use httpmock
	}

	// Create HTTP client with mock authenticator
	client := restyutil.NewHTTPClient(cfg, &mockAuthenticator{})
	httpmock.ActivateNonDefault(client.GetClient())
	defer httpmock.DeactivateAndReset()

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Mock HTTP response
			httpmock.RegisterResponder(
				"GET",
				"http://mock-api/test",
				httpmock.NewJsonResponderOrPanic(tc.statusCode, tc.responseBody),
			)

			// Make request
			resp, err := client.R().Get("/test")

			// Assert errors
			if tc.expectedError != nil {
				require.Error(t, err)

				var spvErr models.SPVError
				if errors.As(err, &spvErr) && tc.expectedSPVError != nil {
					require.Equal(t, *tc.expectedSPVError, spvErr)
				} else {
					require.Contains(t, err.Error(), tc.expectedError.Error())
				}
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
			}
		})
	}
}
