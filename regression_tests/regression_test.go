//go:build regression
// +build regression

package regressiontests

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	wallet "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/config"
	"github.com/bitcoin-sv/spv-wallet-go-client/walletkeys"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/stretchr/testify/require"
)

const (
	minimalFundsPerTransaction = 2

	adminXPriv = "xprv9s21ZrQH143K3CbJXirfrtpLvhT3Vgusdo8coBritQ3rcS7Jy7sxWhatuxG5h2y1Cqj8FKmPp69536gmjYRpfga2MJdsGyBsnB12E19CESK"
	adminXPub  = "xpub661MyMwAqRbcFgfmdkPgE2m5UjHXu9dj124DbaGLSjaqVESTWfCD4VuNmEbVPkbYLCkykwVZvmA8Pbf8884TQr1FgdG2nPoHR8aB36YdDQh"

	errGettingEnvVariables = "failed to get environment variables: %s"
	errGettingSharedConfig = "failed to get shared config: %s"
	errCreatingUser        = "failed to create user: %s"
	errDeletingUserPaymail = "failed to delete user's paymail: %s"
	errSendingFunds        = "failed to send funds: %s"
	errGettingBalance      = "failed to get balance: %s"
	errGettingTransactions = "failed to get transactions: %s"
)

func TestRegression(t *testing.T) {
	ctx := context.Background()
	rtConfig, err := getEnvVariables()
	require.NoError(t, err, fmt.Sprintf(errGettingEnvVariables, err))

	var paymailDomainInstanceOne, paymailDomainInstanceTwo string
	var userOne, userTwo *regressionTestUser

	t.Run("Initialize Shared Configurations", func(t *testing.T) {
		t.Run("Should get sharedConfig for instance one", func(t *testing.T) {
			paymailDomainInstanceOne = getPaymailDomain(ctx, t, adminXPriv, rtConfig.ClientOneURL)
		})

		t.Run("Should get shared config for instance two", func(t *testing.T) {
			paymailDomainInstanceTwo = getPaymailDomain(ctx, t, adminXPriv, rtConfig.ClientTwoURL)
		})
	})

	t.Run("Create Users", func(t *testing.T) {
		t.Run("Should create user for instance one", func(t *testing.T) {
			userName := "instanceOneUser1"
			userOne = createUser(ctx, t, userName, paymailDomainInstanceOne, rtConfig.ClientOneURL, adminXPriv)

		})

		t.Run("Should create user for instance two", func(t *testing.T) {
			userName := "instanceTwoUser1"
			userTwo, err = createUser(ctx, userName, paymailDomainInstanceTwo, rtConfig.ClientTwoURL, adminXPriv)

		})
	})

	defer func() {
		t.Run("Cleanup: Remove Paymails", func(t *testing.T) {
			t.Run("Should remove user's paymail on first instance", func(t *testing.T) {
				if userOne != nil {
					err := removeRegisteredPaymail(ctx, userOne.Paymail, rtConfig.ClientOneURL, adminXPriv)
					require.NoError(t, err, fmt.Sprintf(errDeletingUserPaymail, err))
				}
			})

			t.Run("Should remove user's paymail on second instance", func(t *testing.T) {
				if userTwo != nil {
					err := removeRegisteredPaymail(ctx, userTwo.Paymail, rtConfig.ClientTwoURL, adminXPriv)
					require.NoError(t, err, fmt.Sprintf(errDeletingUserPaymail, err))
				}
			})
		})
	}()

	t.Run("Perform Transactions", func(t *testing.T) {
		t.Run("Send money to instance 1", func(t *testing.T) {
			const amountToSend = 3
			transaction, err := sendFunds(ctx, rtConfig.ClientTwoURL, rtConfig.ClientTwoLeaderXPriv, userOne.Paymail, amountToSend)
			require.NoError(t, err, fmt.Sprintf(errSendingFunds, err))
			require.GreaterOrEqual(t, int64(-1), transaction.OutputValue)

			balance, err := getBalance(ctx, rtConfig.ClientOneURL, userOne.XPriv)
			require.NoError(t, err, fmt.Sprintf(errGettingBalance, err))
			require.GreaterOrEqual(t, balance, 1)

			transactions, err := getTransactions(ctx, rtConfig.ClientOneURL, userOne.XPriv)
			require.NoError(t, err, fmt.Sprintf(errGettingTransactions, err))
			require.GreaterOrEqual(t, len(transactions), 1)
		})

		t.Run("Send money to instance 2", func(t *testing.T) {
			transaction, err := sendFunds(ctx, rtConfig.ClientOneURL, rtConfig.ClientOneLeaderXPriv, userTwo.Paymail, minimalFundsPerTransaction)
			require.NoError(t, err, fmt.Sprintf(errSendingFunds, err))
			require.GreaterOrEqual(t, int64(-1), transaction.OutputValue)

			balance, err := getBalance(ctx, rtConfig.ClientTwoURL, userTwo.XPriv)
			require.NoError(t, err, fmt.Sprintf(errGettingBalance, err))
			require.GreaterOrEqual(t, balance, 1)

			transactions, err := getTransactions(ctx, rtConfig.ClientTwoURL, userTwo.XPriv)
			require.NoError(t, err, fmt.Sprintf(errGettingTransactions, err))
			require.GreaterOrEqual(t, len(transactions), 1)
		})

		t.Run("Send money from instance 1 to instance 2", func(t *testing.T) {
			transaction, err := sendFunds(ctx, rtConfig.ClientOneURL, userOne.XPriv, userTwo.Paymail, minimalFundsPerTransaction)
			require.NoError(t, err, fmt.Sprintf(errSendingFunds, err))
			require.GreaterOrEqual(t, int64(-1), transaction.OutputValue)

			balance, err := getBalance(ctx, rtConfig.ClientTwoURL, userTwo.XPriv)
			require.NoError(t, err, fmt.Sprintf(errGettingBalance, err))
			require.GreaterOrEqual(t, balance, 2)

			transactions, err := getTransactions(ctx, rtConfig.ClientTwoURL, userTwo.XPriv)
			require.NoError(t, err, fmt.Sprintf(errGettingTransactions, err))
			require.GreaterOrEqual(t, len(transactions), 2)

			balance, err = getBalance(ctx, rtConfig.ClientOneURL, userOne.XPriv)
			require.NoError(t, err, fmt.Sprintf(errGettingBalance, err))
			require.GreaterOrEqual(t, balance, 0)

			transactions, err = getTransactions(ctx, rtConfig.ClientOneURL, userOne.XPriv)
			require.NoError(t, err, fmt.Sprintf(errGettingTransactions, err))
			require.GreaterOrEqual(t, len(transactions), 2)
		})
	})
}

const (
	atSign       = "@"
	domainPrefix = "https://"

	ClientOneURLEnvVar         = "CLIENT_ONE_URL"
	ClientTwoURLEnvVar         = "CLIENT_TWO_URL"
	ClientOneLeaderXPrivEnvVar = "CLIENT_ONE_LEADER_XPRIV"
	ClientTwoLeaderXPrivEnvVar = "CLIENT_TWO_LEADER_XPRIV"
)

var (
	explicitHTTPURLRegex      = regexp.MustCompile(`^https?://`)
	errEmptyXPrivEnvVariables = errors.New("missing xpriv variables")
)

type regressionTestUser struct {
	XPriv   string `json:"xpriv"`
	XPub    string `json:"xpub"`
	Paymail string `json:"paymail"`
}

type regressionTestConfig struct {
	ClientOneURL         string
	ClientTwoURL         string
	ClientOneLeaderXPriv string
	ClientTwoLeaderXPriv string
}

// getEnvVariables retrieves the environment variables needed for the regression tests.
func getEnvVariables() (*regressionTestConfig, error) {
	rtConfig := regressionTestConfig{
		ClientOneURL:         os.Getenv(ClientOneURLEnvVar),
		ClientTwoURL:         os.Getenv(ClientTwoURLEnvVar),
		ClientOneLeaderXPriv: os.Getenv(ClientOneLeaderXPrivEnvVar),
		ClientTwoLeaderXPriv: os.Getenv(ClientTwoLeaderXPrivEnvVar),
	}

	if rtConfig.ClientOneLeaderXPriv == "" || rtConfig.ClientTwoLeaderXPriv == "" {
		return nil, errEmptyXPrivEnvVariables
	}
	if rtConfig.ClientOneURL == "" || rtConfig.ClientTwoURL == "" {
		rtConfig.ClientOneURL = "http://localhost:3003"
		rtConfig.ClientTwoURL = "http://localhost:3003"
	}

	rtConfig.ClientOneURL = addPrefixIfNeeded(rtConfig.ClientOneURL)
	rtConfig.ClientTwoURL = addPrefixIfNeeded(rtConfig.ClientTwoURL)

	return &rtConfig, nil
}

// getPaymailDomain retrieves the shared configuration from the SPV Wallet.
func getPaymailDomain(ctx context.Context, t *testing.T, xpriv string, url string) string {
	t.Helper()

	api, err := wallet.NewUserAPIWithXPriv(config.New(config.WithAddr(url)), xpriv)
	if err != nil {
		t.Fatalf("test helper - failed to initialize user API with XPriv: %v", err)
	}

	sharedConfig, err := api.SharedConfig(ctx)
	if err != nil {
		t.Fatalf("test helper - failed to retrieve shared config from User API: %v", err)
	}
	if len(sharedConfig.PaymailDomains) != 1 {
		t.Fatalf("test helper - expected to have single paymail domain, got: %d", len(sharedConfig.PaymailDomains))
	}

	return sharedConfig.PaymailDomains[0]
}

// createUser creates a set of keys and new paymail in the SPV Wallet.
func createUser(ctx context.Context, t *testing.T, paymail string, paymailDomain string, instanceUrl string, adminXPriv string) *regressionTestUser {
	t.Helper()

	keys, err := walletkeys.RandomKeys()
	if err != nil {
		t.Fatalf("test helper - failed to generate random keys: %v", err)
	}

	user := &regressionTestUser{
		XPriv:   keys.XPriv(),
		XPub:    keys.XPub(),
		Paymail: preparePaymail(paymail, paymailDomain),
	}

	cfg := config.New(config.WithAddr(instanceUrl))
	t.Log("test helper - addr: %v", cfg.Addr)

	adminClient, err := wallet.NewAdminAPIWithXPriv(cfg, adminXPriv)
	if err != nil {
		t.Fatalf("test helper - failed to initialize admin API with XPriv: %v", err)
	}

	_, err = adminClient.CreateXPub(ctx, &commands.CreateUserXpub{
		Metadata: map[string]any{"some_metadata": "remove"},
		XPub:     user.XPub,
	})
	if err != nil {
		t.Fatalf("test helper - failed to create XPub: %v", err)
	}

	_, err = adminClient.CreatePaymail(ctx, &commands.CreatePaymail{
		Metadata:   map[string]any{},
		Key:        user.XPub,
		Address:    user.Paymail,
		PublicName: "Regression tests",
		Avatar:     "",
	})
	if err != nil {
		t.Fatalf("test helper - failed to create paymail: %v", err)
	}

	return user
}

// removeRegisteredPaymail soft deletes paymail from the SPV Wallet.
func removeRegisteredPaymail(ctx context.Context, paymail string, instanceUrl string, adminXPriv string) error {
	cfg := config.New(config.WithAddr(instanceUrl))
	adminClient, err := wallet.NewAdminAPIWithXPriv(cfg, adminXPriv)
	if err != nil {
		return err
	}
	err = adminClient.DeletePaymail(ctx, paymail)
	if err != nil {
		return err
	}
	return nil
}

// getBalance retrieves the balance from the SPV Wallet.
func getBalance(ctx context.Context, fromInstance string, fromXPriv string) (int, error) {
	cfg := config.New(config.WithAddr(fromInstance))
	client, err := wallet.NewUserAPIWithXPriv(cfg, fromXPriv)
	if err != nil {
		return -1, err
	}
	xpubInfo, err := client.XPub(ctx)
	if err != nil {
		return -1, err
	}
	return int(xpubInfo.CurrentBalance), nil
}

// getTransactions retrieves the transactions from the SPV Wallet.
func getTransactions(ctx context.Context, fromInstance string, fromXPriv string) ([]*response.Transaction, error) {
	cfg := config.New(config.WithAddr(fromInstance))
	client, err := wallet.NewUserAPIWithXPriv(cfg, fromXPriv)
	if err != nil {
		return nil, err
	}

	page, err := client.Transactions(ctx)
	if err != nil {
		return nil, err
	}
	return page.Content, nil
}

// sendFunds sends funds from one paymail to another.
func sendFunds(ctx context.Context, fromInstance string, fromXPriv string, toPaymail string, howMuch int) (*response.Transaction, error) {
	cfg := config.New(config.WithAddr(fromInstance))
	client, err := wallet.NewUserAPIWithXPriv(cfg, fromXPriv)
	if err != nil {
		return nil, err
	}

	balance, err := getBalance(ctx, fromInstance, fromXPriv)
	if err != nil {
		return nil, err
	}
	if balance < howMuch {
		return nil, fmt.Errorf("insufficient funds: %d", balance)
	}

	transaction, err := client.SendToRecipients(ctx, &commands.SendToRecipients{
		Recipients: []*commands.Recipients{
			{
				To: toPaymail, Satoshis: uint64(howMuch),
			},
		},
		Metadata: map[string]any{"description": "regression-test"},
	})
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

// preparePaymail prepares the paymail address by combining the alias and domain.
func preparePaymail(paymailAlias string, domain string) string {
	if isValidURL(domain) {
		splitedDomain := strings.SplitAfter(domain, "//")
		domain = splitedDomain[1]
	}
	url := paymailAlias + atSign + domain
	return url
}

// addPrefixIfNeeded adds the HTTPS prefix to the URL if it is missing.
func addPrefixIfNeeded(url string) string {
	if !isValidURL(url) {
		return domainPrefix + url
	}
	return url
}

// isValidURL validates the URL if it has http or https prefix.
func isValidURL(rawURL string) bool {
	return explicitHTTPURLRegex.MatchString(rawURL)
}
