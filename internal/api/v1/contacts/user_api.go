package contacts

import (
	"context"
	"fmt"
	"net/url"

	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/contacts/builders"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/errutil"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/querybuilders"
	"github.com/bitcoin-sv/spv-wallet-go-client/queries"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/go-resty/resty/v2"
)

const (
	userAPIRoute = "api/v1/contacts"
	userAPI      = "User Contacts API"
)

type UserAPI struct {
	url        *url.URL
	httpClient *resty.Client
}

func (u *UserAPI) Contacts(ctx context.Context, opts ...queries.ContactQueryOption) (*queries.UserContactsPage, error) {
	var query queries.ContactQuery
	for _, o := range opts {
		o(&query)
	}

	queryBuilder := querybuilders.NewQueryBuilder(
		querybuilders.WithMetadataFilter(query.Metadata),
		querybuilders.WithPageFilter(query.PageFilter),
		querybuilders.WithFilterQueryBuilder(&builders.UserContactFilterQueryBuilder{
			ContactFilter: query.ContactFilter,
			ModelFilterBuilder: querybuilders.ModelFilterBuilder{
				ModelFilter: query.ContactFilter.ModelFilter,
			},
		}),
	)
	params, err := queryBuilder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build user contacts query params: %w", err)
	}

	var result queries.UserContactsPage
	_, err = u.httpClient.
		R().
		SetContext(ctx).
		SetResult(&result).
		SetQueryParams(params.ParseToMap()).
		Get(u.url.String())
	if err != nil {
		return nil, fmt.Errorf("HTTP response failure: %w", err)
	}

	return &result, nil
}

func (u *UserAPI) ContactWithPaymail(ctx context.Context, paymail string) (*response.Contact, error) {
	var result response.Contact
	_, err := u.httpClient.
		R().
		SetContext(ctx).
		SetResult(&result).
		Get(u.url.JoinPath(paymail).String())
	if err != nil {
		return nil, fmt.Errorf("HTTP response failure: %w", err)
	}

	return &result, nil
}

func (u *UserAPI) UpsertContact(ctx context.Context, cmd commands.UpsertContact) (*response.Contact, error) {
	var result response.CreateContactResponse
	_, err := u.httpClient.
		R().
		SetBody(cmd).
		SetContext(ctx).
		SetResult(&result).
		Put(u.url.JoinPath(cmd.Paymail).String())
	if err != nil {
		return nil, fmt.Errorf("HTTP response failure: %w", err)
	}

	return &response.Contact{
		Model:    result.Contact.Model,
		ID:       result.Contact.ID,
		FullName: result.Contact.FullName,
		Paymail:  result.Contact.Paymail,
		PubKey:   result.Contact.PubKey,
		Status:   result.Contact.Status,
	}, nil
}

func (u *UserAPI) RemoveContact(ctx context.Context, paymail string) error {
	_, err := u.httpClient.
		R().
		SetContext(ctx).
		Delete(u.url.JoinPath(paymail).String())
	if err != nil {
		return fmt.Errorf("HTTP response failure: %w", err)
	}

	return nil
}

func (u *UserAPI) ConfirmContact(ctx context.Context, paymail string) error {
	_, err := u.httpClient.
		R().
		SetContext(ctx).
		Post(u.url.JoinPath(paymail, "confirmation").String())
	if err != nil {
		return fmt.Errorf("HTTP response failure: %w", err)
	}

	return nil
}

func (u *UserAPI) UnconfirmContact(ctx context.Context, paymail string) error {
	_, err := u.httpClient.
		R().
		SetContext(ctx).
		Delete(u.url.JoinPath(paymail, "confirmation").String())
	if err != nil {
		return fmt.Errorf("HTTP response failure: %w", err)
	}

	return nil
}

func NewUserAPI(url *url.URL, httpClient *resty.Client) *UserAPI {
	return &UserAPI{
		url:        url.JoinPath(userAPIRoute),
		httpClient: httpClient,
	}
}

func UserAPIErrorFormatter(action string, err error) *errutil.HTTPErrorFormatter {
	return &errutil.HTTPErrorFormatter{
		Action: action,
		API:    userAPI,
		Err:    err,
	}
}
