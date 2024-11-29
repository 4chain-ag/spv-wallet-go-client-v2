package contacts_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/errors"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/admin/contacts/contactstest"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/spvwallettest"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestContactsAPI_Contacts(t *testing.T) {
	tests := map[string]struct {
		responder        httpmock.Responder
		expectedResponse *queries.UserContactsPage
		expectedErr      error
	}{
		"HTTP GET /api/v1/admin/contacts response: 200": {
			expectedResponse: contactstest.ExpectedUserContactsPage(t),
			responder:        httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("contactstest/get_contacts_200.json")),
		},
		"HTTP GET /api/v1/admin/contacts response: 400": {
			expectedErr: contactstest.NewBadRequestSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, contactstest.NewBadRequestSPVError()),
		},
		"HTTP GET /api/v1/admin/contacts str response: 500": {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	URL := spvwallettest.TestAPIAddr + "/api/v1/admin/contacts"
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			wallet, transport := spvwallettest.GivenSPVAdminAPI(t)
			transport.RegisterResponder(http.MethodGet, URL, tc.responder)

			// then:
			got, err := wallet.Contacts(context.Background(), queries.ContactQueryWithPageFilter(filter.Page{Size: 1}))
			require.ErrorIs(t, err, tc.expectedErr)
			require.EqualValues(t, tc.expectedResponse, got)
		})
	}
}

func TestContactsAPI_ContactUpdate(t *testing.T) {
	ID := "4d570959-dd85-4f53-bad1-18d0671761e9"
	tests := map[string]struct {
		responder        httpmock.Responder
		expectedResponse *response.Contact
		expectedErr      error
	}{
		fmt.Sprintf("HTTP PUT /api/v1/admin/contacts/%s response: 200", ID): {
			expectedResponse: contactstest.ExpectedUpdatedUserContact(t),
			responder:        httpmock.NewJsonResponderOrPanic(http.StatusOK, httpmock.File("contactstest/put_contact_update_200.json")),
		},
		fmt.Sprintf("HTTP PUT /api/v1/admin/contacts/%s response: 400", ID): {
			expectedErr: contactstest.NewBadRequestSPVError(),
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, contactstest.NewBadRequestSPVError()),
		},
		fmt.Sprintf("HTTP PUT /api/v1/admin/contacts/%s response: 500", ID): {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	URL := spvwallettest.TestAPIAddr + "/api/v1/admin/contacts/" + ID
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			wallet, transport := spvwallettest.GivenSPVAdminAPI(t)
			transport.RegisterResponder(http.MethodPut, URL, tc.responder)

			// then:
			got, err := wallet.ContactUpdate(context.Background(), &commands.UpdateContact{
				ID:       "4d570959-dd85-4f53-bad1-18d0671761e9",
				FullName: "John Doe Williams",
				Metadata: map[string]any{"phoneNumber": "123456789"},
			})
			require.ErrorIs(t, err, tc.expectedErr)
			require.EqualValues(t, tc.expectedResponse, got)
		})
	}
}

func TestContactsAPI_RemoveContact(t *testing.T) {
	ID := "4d570959-dd85-4f53-bad1-18d0671761e9"
	tests := map[string]struct {
		responder   httpmock.Responder
		statusCode  int
		expectedErr error
	}{
		fmt.Sprintf("HTTP DELETE /api/v1/admin/contacts/%s response: 200", ID): {
			statusCode: http.StatusOK,
			responder:  httpmock.NewStringResponder(http.StatusOK, http.StatusText(http.StatusOK)),
		},
		fmt.Sprintf("HTTP DELETE/api/v1/admin/contacts/%s response: 400", ID): {
			expectedErr: contactstest.NewBadRequestSPVError(),
			statusCode:  http.StatusOK,
			responder:   httpmock.NewJsonResponderOrPanic(http.StatusBadRequest, contactstest.NewBadRequestSPVError()),
		},
		fmt.Sprintf("HTTP DELETE /api/v1/admin/contacts/%s str response: 500", ID): {
			expectedErr: errors.ErrUnrecognizedAPIResponse,
			statusCode:  http.StatusInternalServerError,
			responder:   httpmock.NewStringResponder(http.StatusInternalServerError, "unexpected internal server failure"),
		},
	}

	url := spvwallettest.TestAPIAddr + "/api/v1/admin/contacts/" + ID
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			wallet, transport := spvwallettest.GivenSPVAdminAPI(t)
			transport.RegisterResponder(http.MethodDelete, url, tc.responder)

			// then:
			err := wallet.RemoveContact(context.Background(), ID)
			require.ErrorIs(t, err, tc.expectedErr)
		})
	}
}
