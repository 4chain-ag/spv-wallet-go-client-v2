package users_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	client "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/users/userstest"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/testfixtures"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestUsersAPI_UpdateXPubMetadata(t *testing.T) {
	tests := map[string]struct {
		code             int
		responder        httpmock.Responder
		statusCode       int
		expectedResponse *response.Xpub
		expectedErr      error
	}{
		"HTTP PATCH /api/v1/users/current response: 200": {
			expectedResponse: userstest.ExpectedUpdatedXPubMetadata(t),
			code:             http.StatusOK,
			responder:        httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("userstest/updated_xpub_metadata_200.json")),
		},
		"HTTP PATCH /api/v1/users/current response: 400": {
			expectedErr: models.SPVError{
				Message:    http.StatusText(http.StatusBadRequest),
				StatusCode: http.StatusBadRequest,
				Code:       "invalid-data-format",
			},
			statusCode: http.StatusOK,
			responder:  httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, userstest.NewBadRequestSPVError()),
		},
		"HTTP PATCH /api/v1/users/current str response: 500": {
			expectedErr: client.ErrUnrecognizedAPIResponse,
			statusCode:  http.StatusInternalServerError,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	URL := testfixtures.TestAPIAddr + "/api/v1/users/current"
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			wallet, transport := testfixtures.GivenSPVWalletClient(t)
			transport.RegisterResponder(http.MethodPatch, URL, tc.responder)

			// then:
			got, err := wallet.UpdateXPubMetadata(context.Background(), &commands.UpdateXPubMetadata{
				Metadata: map[string]any{
					"example_key": "example_value",
				},
			})
			require.ErrorIs(t, err, tc.expectedErr)
			require.EqualValues(t, tc.expectedResponse, got)
		})
	}
}

func TestUsersAPI_XPub(t *testing.T) {
	tests := map[string]struct {
		code             int
		responder        httpmock.Responder
		statusCode       int
		expectedResponse *response.Xpub
		expectedErr      error
	}{
		"HTTP GET /api/v1/users/current/ response: 200": {
			expectedResponse: userstest.ExpectedUserXPub(t),
			code:             http.StatusOK,
			responder:        httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("userstest/retrieved_xpub_200.json")),
		},
		"HTTP GET /api/v1/users/current response: 400": {
			expectedErr: models.SPVError{
				Message:    http.StatusText(http.StatusBadRequest),
				StatusCode: http.StatusBadRequest,
				Code:       "invalid-data-format",
			},
			statusCode: http.StatusOK,
			responder:  httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, userstest.NewBadRequestSPVError()),
		},
		"HTTP GET /api/v1/users/current str response: 500": {
			expectedErr: client.ErrUnrecognizedAPIResponse,
			statusCode:  http.StatusInternalServerError,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	URL := testfixtures.TestAPIAddr + "/api/v1/users/current"
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			wallet, transport := testfixtures.GivenSPVWalletClient(t)
			transport.RegisterResponder(http.MethodGet, URL, tc.responder)

			// then:
			got, err := wallet.XPub(context.Background())
			require.ErrorIs(t, err, tc.expectedErr)
			require.EqualValues(t, tc.expectedResponse, got)
		})
	}
}

func TestUsersAPI_GenerateAccessKey(t *testing.T) {
	tests := map[string]struct {
		code             int
		responder        httpmock.Responder
		statusCode       int
		expectedResponse *response.AccessKey
		expectedErr      error
	}{
		"HTTP POST /api/v1/users/current/keys response: 200": {
			expectedResponse: userstest.ExpectedCreatedAccessKey(t),
			code:             http.StatusOK,
			responder:        httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("userstest/created_access_key_200.json")),
		},
		"HTTP POST /api/v1/users/current/keys response: 400": {
			expectedErr: models.SPVError{
				Message:    http.StatusText(http.StatusBadRequest),
				StatusCode: http.StatusBadRequest,
				Code:       "invalid-data-format",
			},
			statusCode: http.StatusOK,
			responder:  httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, userstest.NewBadRequestSPVError()),
		},
		"HTTP POST /api/v1/users/current/keys str response: 500": {
			expectedErr: client.ErrUnrecognizedAPIResponse,
			statusCode:  http.StatusInternalServerError,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	URL := testfixtures.TestAPIAddr + "/api/v1/users/current/keys"
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			wallet, transport := testfixtures.GivenSPVWalletClient(t)
			transport.RegisterResponder(http.MethodPost, URL, tc.responder)

			// then:
			got, err := wallet.GenerateAccessKey(context.Background(), &commands.GenerateAccessKey{
				Metadata: map[string]any{
					"example_key": "example_value",
				},
			})
			require.ErrorIs(t, err, tc.expectedErr)
			require.EqualValues(t, tc.expectedResponse, got)
		})
	}
}

func TestUsersAPI_AccessKey(t *testing.T) {
	ID := "1fb70cc2-e9d9-41a3-842e-f71cc58d9787"
	tests := map[string]struct {
		code             int
		responder        httpmock.Responder
		statusCode       int
		expectedResponse *response.AccessKey
		expectedErr      error
	}{
		fmt.Sprintf("HTTP GET /api/v1/users/current/keys/%s response: 200", ID): {
			expectedResponse: userstest.ExpectedRertrivedAccessKey(t),
			code:             http.StatusOK,
			responder:        httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("userstest/retrieved_access_key_200.json")),
		},
		fmt.Sprintf("HTTP GET /api/v1/users/current/keys/%s response: 400", ID): {
			expectedErr: models.SPVError{
				Message:    http.StatusText(http.StatusBadRequest),
				StatusCode: http.StatusBadRequest,
				Code:       "invalid-data-format",
			},
			statusCode: http.StatusOK,
			responder:  httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, userstest.NewBadRequestSPVError()),
		},
		fmt.Sprintf("HTTP GET /api/v1/users/current/keys/%s str response: 500", ID): {
			expectedErr: client.ErrUnrecognizedAPIResponse,
			statusCode:  http.StatusInternalServerError,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	URL := testfixtures.TestAPIAddr + "/api/v1/users/current/keys/" + ID
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			wallet, transport := testfixtures.GivenSPVWalletClient(t)
			transport.RegisterResponder(http.MethodGet, URL, tc.responder)

			// then:
			got, err := wallet.AccessKey(context.Background(), ID)
			require.ErrorIs(t, err, tc.expectedErr)
			require.EqualValues(t, tc.expectedResponse, got)
		})
	}
}

func TestUsersAPI_RevokeAccessKey(t *testing.T) {
	ID := "081743f7-040e-45a3-8c36-dde39001e20d"
	tests := map[string]struct {
		code        int
		responder   httpmock.Responder
		statusCode  int
		expectedErr error
	}{
		fmt.Sprintf("HTTP DELETE /api/v1/users/current/keys/%s response: 200", ID): {
			code:      http.StatusOK,
			responder: httpmock.NewStringResponder(http.StatusOK, http.StatusText(http.StatusOK)),
		},
		fmt.Sprintf("HTTP DELETE /api/v1/users/current/keys/%s response: 400", ID): {
			expectedErr: models.SPVError{
				Message:    http.StatusText(http.StatusBadRequest),
				StatusCode: http.StatusBadRequest,
				Code:       "invalid-data-format",
			},
			statusCode: http.StatusOK,
			responder:  httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, userstest.NewBadRequestSPVError()),
		},
		fmt.Sprintf("HTTP DELETE /api/v1/users/current/keys/%s str response: 500", ID): {
			expectedErr: client.ErrUnrecognizedAPIResponse,
			statusCode:  http.StatusInternalServerError,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	URL := testfixtures.TestAPIAddr + "/api/v1/users/current/keys/" + ID
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			wallet, transport := testfixtures.GivenSPVWalletClient(t)
			transport.RegisterResponder(http.MethodDelete, URL, tc.responder)

			// then:
			err := wallet.RevokeAccessKey(context.Background(), ID)
			require.ErrorIs(t, err, tc.expectedErr)
		})
	}
}
