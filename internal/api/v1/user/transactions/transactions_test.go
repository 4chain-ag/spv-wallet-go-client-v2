package transactions_test

import (
	"context"
	"net/http"
	"testing"

	client "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/transactions/transactionsfixture"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/testfixtures"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestTransactionsAPI_Transactions_APIResponses(t *testing.T) {
	tests := map[string]struct {
		statusCode       int
		expectedResponse []*response.Transaction
		expectedErr      error
		responder        httpmock.Responder
	}{
		"HTTP GET /api/v1/transactions response: 200": {
			statusCode:       http.StatusOK,
			expectedResponse: transactionsfixture.NewTransactions().Content,
			responder:        httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("transactionsfixture/responses/transactions_200_status_code.json")),
		},
		"HTTP GET /api/v1/transactions response: 400": {
			expectedErr: models.SPVError{
				Message:    http.StatusText(http.StatusBadRequest),
				StatusCode: http.StatusBadRequest,
				Code:       "invalid-data-format",
			},
			statusCode: http.StatusOK,
			responder: httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, &models.SPVError{
				Message:    http.StatusText(http.StatusBadRequest),
				StatusCode: http.StatusBadRequest,
				Code:       "invalid-data-format",
			}),
		},
		"HTTP GET /api/v1/transactions str response: 500": {
			expectedErr: client.ErrUnrecognizedAPIResponse,
			statusCode:  http.StatusInternalServerError,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	URL := testfixtures.TestAPIAddr + "/api/v1/transactions"
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			wallet, transport := testfixtures.GivenSPVWalletClient(t)
			transport.RegisterResponder(http.MethodGet, URL, tc.responder)

			// then:
			got, err := wallet.Transactions(context.Background())
			require.ErrorIs(t, err, tc.expectedErr)
			require.EqualValues(t, tc.expectedResponse, got)
		})
	}
}
