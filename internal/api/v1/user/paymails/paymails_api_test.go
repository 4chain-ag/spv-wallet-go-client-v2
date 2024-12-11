package paymails_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/errors"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/paymails/paymailstest"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/testutils"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

const paymailsURL = "/api/v1/paymails"

func TestPaymailsAPI_Paymails(t *testing.T) {
	tests := map[string]struct {
		responder        httpmock.Responder
		expectedResponse *queries.PaymailAddressPage
		expectedErr      error
	}{
		"HTTP GET /api/v1/paymails response: 200": {
			expectedResponse: paymailstest.ExpectedPaymailsPage(t),
			responder:        testutils.NewJSONFileResponderWithStatusOK("paymailstest/get_paymails_200.json"),
		},
		"HTTP GET /api/v1/paymails response: 400": {
			expectedErr: testutils.NewBadRequestSPVError(),
			responder:   testutils.NewBadRequestSPVErrorResponder(),
		},
		"HTTP GET /api/v1/paymails response: 500": {
			expectedErr: testutils.NewInternalServerSPVError(),
			responder:   testutils.NewInternalServerSPVErrorResponder(),
		},
		"HTTP GET /api/v1/paymails str response: 500": {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   testutils.NewInternalServerSPVErrorStringResponder("unexpected internal server failure"),
		},
	}

	url := testutils.FullAPIURL(t, paymailsURL)
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			wallet, transport := testutils.GivenSPVUserAPI(t)
			transport.RegisterResponder(http.MethodGet, url, tc.responder)

			// then:
			got, err := wallet.Paymails(context.Background(), queries.PaymailQueryWithPageFilter(filter.Page{Size: 1}))
			require.ErrorIs(t, err, tc.expectedErr)
			require.EqualValues(t, tc.expectedResponse, got)
		})
	}
}
