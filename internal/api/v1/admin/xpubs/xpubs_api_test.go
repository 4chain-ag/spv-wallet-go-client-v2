package xpubs_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/errors"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/admin/xpubs/xpubstest"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/spvwallettest"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestXPubsAPI_CreateXPub(t *testing.T) {
	tests := map[string]struct {
		responder        httpmock.Responder
		expectedResponse *response.Xpub
		expectedErr      error
	}{
		"HTTP POST /api/v1/admin/users response: 201": {
			expectedResponse: xpubstest.ExpectedXPub(t),
			responder:        httpmock.NewJsonResponderOrPanic(http.StatusCreated, httpmock.File("xpubstest/post_xpub_201.json")),
		},
		"HTTP POST /api/v1/admin/users response: 400": {
			expectedErr: xpubstest.NewBadRequestSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, xpubstest.NewBadRequestSPVError()),
		},
		"HTTP POST /api/v1/admin/users str response: 500": {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	URL := spvwallettest.TestAPIAddr + "/api/v1/admin/users"
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			wallet, transport := spvwallettest.GivenSPVAdminAPI(t)
			transport.RegisterResponder(http.MethodPost, URL, tc.responder)

			// then:
			got, err := wallet.CreateXPub(context.Background(), &commands.CreateUserXpub{
				Metadata: map[string]any{},
				XPub:     "",
			})
			require.ErrorIs(t, err, tc.expectedErr)
			require.EqualValues(t, tc.expectedResponse, got)
		})
	}
}

func TestXPubsAPI_XPubs(t *testing.T) {
	tests := map[string]struct {
		responder        httpmock.Responder
		expectedResponse *queries.XPubPage
		expectedErr      error
	}{
		"HTTP GET /api/v1/admin/users response: 200": {
			expectedResponse: xpubstest.ExpectedXPubsPage(t),
			responder:        httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("xpubstest/get_xpubs_200.json")),
		},
		"HTTP GET /api/v1/admin/users response: 400": {
			expectedErr: xpubstest.NewBadRequestSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, xpubstest.NewBadRequestSPVError()),
		},
		"HTTP GET /api/v1/admin/users str response: 500": {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	URL := spvwallettest.TestAPIAddr + "/api/v1/admin/users"
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			wallet, transport := spvwallettest.GivenSPVAdminAPI(t)
			transport.RegisterResponder(http.MethodGet, URL, tc.responder)

			// then:
			got, err := wallet.XPubs(context.Background())
			require.ErrorIs(t, err, tc.expectedErr)
			require.EqualValues(t, tc.expectedResponse, got)
		})
	}
}
