package merkleroots_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/errors"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/merkleroots/merklerootstest"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/spvwallettest"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestMerkleRootsAPI_MerkleRoots(t *testing.T) {
	tests := map[string]struct {
		responder        httpmock.Responder
		statusCode       int
		expectedResponse *queries.MerkleRootPage
		expectedErr      error
	}{
		"HTTP GET /api/v1/merkleroots response: 200": {
			statusCode:       http.StatusOK,
			expectedResponse: merklerootstest.ExpectedMerkleRootsPage(),
			responder:        httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("merklerootstest/get_merkleroots_200.json")),
		},
		"HTTP GET /api/v1/merkleroots response: 400": {
			expectedErr: models.SPVError{
				Message:    http.StatusText(http.StatusBadRequest),
				StatusCode: http.StatusBadRequest,
				Code:       "invalid-data-format",
			},
			statusCode: http.StatusOK,
			responder:  httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, merklerootstest.NewBadRequestSPVError()),
		},
		"HTTP GET /api/v1/merkleroots str response: 500": {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			statusCode:  http.StatusInternalServerError,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	url := spvwallettest.TestAPIAddr + "/api/v1/merkleroots"
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			spvWalletClient, transport := spvwallettest.GivenSPVWalletClient(t)
			transport.RegisterResponder(http.MethodGet, url, tc.responder)

			// then:
			got, err := spvWalletClient.MerkleRoots(context.Background())
			require.ErrorIs(t, err, tc.expectedErr)
			require.EqualValues(t, tc.expectedResponse, got)
		})
	}
}
