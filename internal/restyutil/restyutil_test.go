package restyutil_test

import (
	"net/http"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/restyutil"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestResponseAdapter_HandleErr(t *testing.T) {
	httpmock.Activate()
	t.Cleanup(httpmock.DeactivateAndReset)

	tests := map[string]struct {
		expectedErr error
		str         string
		statusCode  int
		responder   httpmock.Responder
	}{
		"HTTP GET: response status code 200": {
			statusCode: http.StatusOK,
		},
		"HTTP GET: response status code 400": {
			statusCode: http.StatusBadRequest,
			expectedErr: models.SPVError{
				Code:       "missing-input-paramteres",
				Message:    http.StatusText(http.StatusBadRequest),
				StatusCode: http.StatusBadRequest,
			},
		},
		"HTTP GET: response status code 500": {
			statusCode: http.StatusBadRequest,
			expectedErr: models.SPVError{
				Code:       "internal-server-error",
				Message:    http.StatusText(http.StatusInternalServerError),
				StatusCode: http.StatusInternalServerError,
			},
		},
		"HTTP GET: API str err response": {
			str:         string("error"),
			statusCode:  http.StatusBadRequest,
			expectedErr: restyutil.ErrUnrecognizedAPIResponse,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			transport := httpmock.NewMockTransport()
			transport.RegisterResponder(http.MethodGet, "/", httpResponderTestHelper(t, tc.expectedErr, tc.str, tc.statusCode))

			cli := resty.
				New().
				SetTransport(transport).
				SetError(&models.SPVError{})

			response := restyutil.ResponseAdapter(func() (*resty.Response, error) { return cli.R().Get("/") })
			err := response.HandleErr()
			require.ErrorIs(t, err, tc.expectedErr)
		})
	}
}

func httpResponderTestHelper(t *testing.T, err error, s string, code int) httpmock.Responder {
	t.Helper()
	return func(r *http.Request) (*http.Response, error) {
		switch {
		case err != nil:
			resp, err := httpmock.NewJsonResponse(code, err)
			if err != nil {
				t.Fatalf("failed to create JSON error response: %s", err)
			}
			return resp, nil
		case s != "":
			return httpmock.NewStringResponse(code, s), nil
		default:
			return httpmock.NewStringResponse(code, "{}"), nil
		}
	}
}
