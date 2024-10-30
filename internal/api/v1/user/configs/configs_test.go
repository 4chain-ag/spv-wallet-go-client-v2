package configs_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	client "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestConfigsAPI_SharedConfig(t *testing.T) {
	httpmock.Activate()
	t.Cleanup(httpmock.DeactivateAndReset)

	tests := map[string]struct {
		expectedErr error
		expectedCfg *response.SharedConfig
	}{
		"shared config: HTTP GET /api/v1/configs/shared response: 200": {
			expectedCfg: &response.SharedConfig{
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
			expectedErr: models.SPVError{
				Message:    http.StatusText(http.StatusInternalServerError),
				StatusCode: http.StatusInternalServerError,
				Code:       models.UnknownErrorCode,
			},
		},
	}

	transport := httpmock.NewMockTransport()
	cfg := client.Config{
		Addr:      "http://localhost:3003",
		Timeout:   time.Minute,
		Transport: transport,
	}
	URL := cfg.Addr + "/api/v1/configs/shared"
	SPV := newClientTestHelper(t, cfg)

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			transport.RegisterResponder(http.MethodGet, URL, httpResponderTestHelper(t, tc.expectedCfg, tc.expectedErr))
			got, err := SPV.SharedConfig(context.Background())
			require.ErrorIs(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedCfg, got)
		})
	}
}

func httpResponderTestHelper(t *testing.T, expectedRes *response.SharedConfig, expectedErr error) httpmock.Responder {
	t.Helper()
	return func(r *http.Request) (*http.Response, error) {
		if expectedErr != nil {
			resp, err := httpmock.NewJsonResponse(http.StatusBadRequest, expectedErr)
			if err != nil {
				t.Fatalf("failed to create JSON error response: %s", err)
			}
			return resp, nil
		}

		resp, err := httpmock.NewJsonResponse(http.StatusOK, expectedRes)
		if err != nil {
			t.Fatalf("failed to create JSON response: %s", err)
		}
		return resp, nil
	}
}

func newClientTestHelper(t *testing.T, cfg client.Config) *client.Client {
	t.Helper()
	xPriv := "xprv9s21ZrQH143K3fqNnUmXmgfT9ToMtiq5cuKsVBG4E5UqVh4psHDY2XKsEfZKuV4FSZcPS9CYgEQiLUpW2xmHqHFyp23SvTkTCE153cCdwaj"
	spv, err := client.NewWithXPriv(cfg, xPriv)
	if err != nil {
		t.Fatalf("failed to initialize spv wallet client with xpriv: %s", err)
	}
	return spv
}
