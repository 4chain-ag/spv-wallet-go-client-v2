package contacts_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"

	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/contacts/contactstest"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/spvwallettest"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
)

func TestContactsAPI_Contacts(t *testing.T) {
	tests := map[string]struct {
		responder        httpmock.Responder
		expectedResponse *queries.UserContactsPage
		expectedErr      error
	}{
		"HTTP GET /api/v1/contacts response: 200": {
			expectedResponse: contactstest.ExpectedUserContactsPage(t),
			responder:        httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("contactstest/get_contacts_200.json")),
		},
		"HTTP GET /api/v1/contacts response: 400": {
			expectedErr: contactstest.NewBadRequestSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, contactstest.NewBadRequestSPVError()),
		},
		"HTTP GET /api/v1/contacts str response: 500": {
			expectedErr: models.SPVError{
				Message:    http.StatusText(http.StatusInternalServerError),
				StatusCode: http.StatusInternalServerError,
				Code:       models.UnknownErrorCode,
			},
			responder: httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, models.SPVError{
				Message:    http.StatusText(http.StatusInternalServerError),
				StatusCode: http.StatusInternalServerError,
				Code:       models.UnknownErrorCode,
			}),
		},
	}

	url := spvwallettest.TestAPIAddr + "/api/v1/contacts"
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			wallet, transport := spvwallettest.GivenSPVUserAPI(t)
			transport.RegisterResponder(http.MethodGet, url, tc.responder)

			// when:
			got, err := wallet.Contacts(context.Background())

			// then:
			require.ErrorIs(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedResponse, got)
		})
	}
}

func TestContactsAPI_ContactWithPaymail(t *testing.T) {
	paymail := "john.doe.test5@john.doe.test.4chain.space"
	tests := map[string]struct {
		responder        httpmock.Responder
		expectedResponse *response.Contact
		expectedErr      error
	}{
		fmt.Sprintf("HTTP GET /api/v1/contacts/%s response: 200", paymail): {
			expectedResponse: contactstest.ExpectedContactWithWithPaymail(t),
			responder:        httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("contactstest/get_contact_paymail_200.json")),
		},
		fmt.Sprintf("HTTP GET /api/v1/contacts/%s response: 400", paymail): {
			expectedErr: contactstest.NewBadRequestSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, contactstest.NewBadRequestSPVError()),
		},
		fmt.Sprintf("HTTP GET /api/v1/contacts/%s str response: 500", paymail): {
			expectedErr: models.SPVError{
				Message:    http.StatusText(http.StatusInternalServerError),
				StatusCode: http.StatusInternalServerError,
				Code:       models.UnknownErrorCode,
			},
			responder: httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, models.SPVError{
				Message:    http.StatusText(http.StatusInternalServerError),
				StatusCode: http.StatusInternalServerError,
				Code:       models.UnknownErrorCode,
			}),
		},
	}

	url := spvwallettest.TestAPIAddr + "/api/v1/contacts/" + paymail
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			wallet, transport := spvwallettest.GivenSPVUserAPI(t)
			transport.RegisterResponder(http.MethodGet, url, tc.responder)

			// when:
			got, err := wallet.ContactWithPaymail(context.Background(), paymail)

			// then:
			require.ErrorIs(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedResponse, got)
		})
	}
}

func TestContactsAPI_UpsertContact(t *testing.T) {
	paymail := "john.doe.test@john.doe.test.4chain.space"
	tests := map[string]struct {
		responder        httpmock.Responder
		expectedResponse *response.Contact
		expectedErr      error
	}{
		fmt.Sprintf("HTTP PUT /api/v1/contacts/%s response: 200", paymail): {
			expectedResponse: contactstest.ExpectedUpsertContact(t),
			responder:        httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("contactstest/put_contact_upsert_200.json")),
		},
		fmt.Sprintf("HTTP PUT /api/v1/contacts/%s response: 400", paymail): {
			expectedErr: contactstest.NewBadRequestSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, contactstest.NewBadRequestSPVError()),
		},
		fmt.Sprintf("HTTP PUT /api/v1/contacts/%s str response: 500", paymail): {
			expectedErr: models.SPVError{
				Message:    http.StatusText(http.StatusInternalServerError),
				StatusCode: http.StatusInternalServerError,
				Code:       models.UnknownErrorCode,
			},
			responder: httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, models.SPVError{
				Message:    http.StatusText(http.StatusInternalServerError),
				StatusCode: http.StatusInternalServerError,
				Code:       models.UnknownErrorCode,
			}),
		},
	}

	url := spvwallettest.TestAPIAddr + "/api/v1/contacts/" + paymail
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			wallet, transport := spvwallettest.GivenSPVUserAPI(t)
			transport.RegisterResponder(http.MethodPut, url, tc.responder)

			// when:
			got, err := wallet.UpsertContact(context.Background(), commands.UpsertContact{
				FullName: "John Doe",
				Metadata: map[string]any{"example_key": "example_val"},
				Paymail:  paymail,
			})

			// then:
			require.ErrorIs(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedResponse, got)
		})
	}
}

func TestContactsAPI_RemoveContact(t *testing.T) {
	paymail := "john.doe.test@john.doe.test.4chain.space"
	tests := map[string]struct {
		responder   httpmock.Responder
		expectedErr error
	}{
		fmt.Sprintf("HTTP DELETE /api/v1/contacts/%s response: 200", paymail): {
			responder: httpmock.NewStringResponder(http.StatusOK, http.StatusText(http.StatusOK)),
		},
		fmt.Sprintf("HTTP DELETE /api/v1/contacts/%s response: 400", paymail): {
			expectedErr: contactstest.NewBadRequestSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, contactstest.NewBadRequestSPVError()),
		},
		fmt.Sprintf("HTTP DELETE /api/v1/contacts/%s str response: 500", paymail): {
			expectedErr: models.SPVError{
				Message:    http.StatusText(http.StatusInternalServerError),
				StatusCode: http.StatusInternalServerError,
				Code:       models.UnknownErrorCode,
			},
			responder: httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, models.SPVError{
				Message:    http.StatusText(http.StatusInternalServerError),
				StatusCode: http.StatusInternalServerError,
				Code:       models.UnknownErrorCode,
			}),
		},
	}

	url := spvwallettest.TestAPIAddr + "/api/v1/contacts/" + paymail
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			wallet, transport := spvwallettest.GivenSPVUserAPI(t)
			transport.RegisterResponder(http.MethodDelete, url, tc.responder)

			// when:
			err := wallet.RemoveContact(context.Background(), paymail)

			// then:
			require.ErrorIs(t, err, tc.expectedErr)
		})
	}
}

func TestContactsAPI_ConfirmContact(t *testing.T) {
	contact := &models.Contact{
		Paymail: "alice@example.com",
		PubKey:  spvwallettest.MockPKI(t, spvwallettest.UserXPub),
	}

	tests := map[string]struct {
		responder   httpmock.Responder
		expectedErr error
	}{
		fmt.Sprintf("HTTP POST /api/v1/contacts/%s/confirmation response: 200", contact.Paymail): {
			responder: httpmock.NewStringResponder(http.StatusOK, http.StatusText(http.StatusOK)),
		},
		fmt.Sprintf("HTTP POST /api/v1/contacts/%s/confirmation response: 400", contact.Paymail): {
			expectedErr: contactstest.NewBadRequestSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, contactstest.NewBadRequestSPVError()),
		},
		fmt.Sprintf("HTTP POST /api/v1/contacts/%s/confirmation str response: 500", contact.Paymail): {
			expectedErr: models.SPVError{
				Message:    http.StatusText(http.StatusInternalServerError),
				StatusCode: http.StatusInternalServerError,
				Code:       models.UnknownErrorCode,
			},
			responder: httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, models.SPVError{
				Message:    http.StatusText(http.StatusInternalServerError),
				StatusCode: http.StatusInternalServerError,
				Code:       models.UnknownErrorCode,
			}),
		},
	}

	url := spvwallettest.TestAPIAddr + "/api/v1/contacts/" + contact.Paymail + "/confirmation"
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			const period = 3600
			const digits = 6

			wrappedTransport := spvwallettest.NewTransportWrapper()
			aliceClient, _ := spvwallettest.GivenSPVWalletClientWithTransport(t, wrappedTransport)
			wrappedTransport.RegisterResponder(http.MethodPost, url, tc.responder)

			// when:
			passcode, err := aliceClient.GenerateTotpForContact(contact, period, digits)

			// then:
			require.NoError(t, err)
			require.NotEmpty(t, passcode)

			err = aliceClient.ConfirmContact(context.Background(), contact, passcode, contact.Paymail, period, digits)
			require.ErrorIs(t, err, tc.expectedErr)
		})
	}
}

func TestContactsAPI_UnconfirmContact(t *testing.T) {
	paymail := "john.doe.test@john.doe.test.4chain.space"
	tests := map[string]struct {
		responder   httpmock.Responder
		expectedErr error
	}{
		fmt.Sprintf("HTTP DELETE /api/v1/contacts/%s/confirmation response: 200", paymail): {
			responder: httpmock.NewStringResponder(http.StatusOK, http.StatusText(http.StatusOK)),
		},
		fmt.Sprintf("HTTP DELETE /api/v1/contacts/%s/confirmation response: 400", paymail): {
			expectedErr: contactstest.NewBadRequestSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, contactstest.NewBadRequestSPVError()),
		},
		fmt.Sprintf("HTTP DELETE /api/v1/contacts/%s/confirmation str response: 500", paymail): {
			expectedErr: models.SPVError{
				Message:    http.StatusText(http.StatusInternalServerError),
				StatusCode: http.StatusInternalServerError,
				Code:       models.UnknownErrorCode,
			},
			responder: httpmock.NewJsonResponderOrPanic(http.StatusInternalServerError, models.SPVError{
				Message:    http.StatusText(http.StatusInternalServerError),
				StatusCode: http.StatusInternalServerError,
				Code:       models.UnknownErrorCode,
			}),
		},
	}

	url := spvwallettest.TestAPIAddr + "/api/v1/contacts/" + paymail + "/confirmation"
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			wallet, transport := spvwallettest.GivenSPVUserAPI(t)
			transport.RegisterResponder(http.MethodDelete, url, tc.responder)

			// when:
			err := wallet.UnconfirmContact(context.Background(), paymail)

			// then:
			require.ErrorIs(t, err, tc.expectedErr)
		})
	}
}
