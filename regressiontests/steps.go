package regressiontests

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	wallet "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/config"
	"github.com/bitcoin-sv/spv-wallet-go-client/walletkeys"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

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
func getPaymailDomain(ctx context.Context, xpriv string, clientUrl string) (string, error) {
	cfg := config.NewDefaultConfig(clientUrl)
	wc, err := wallet.NewUserAPIWithXPriv(cfg, xpriv)
	if err != nil {
		return "", err
	}
	sharedConfig, err := wc.SharedConfig(ctx)
	if err != nil {
		return "", err
	}
	if len(sharedConfig.PaymailDomains) != 1 {
		return "", fmt.Errorf("expected 1 paymail domain, got %d", len(sharedConfig.PaymailDomains))
	}
	return sharedConfig.PaymailDomains[0], nil
}

// createUser creates a set of keys and new paymail in the SPV Wallet.
func createUser(ctx context.Context, paymail string, paymailDomain string, instanceUrl string, adminXPriv string) (*regressionTestUser, error) {
	keys, err := walletkeys.RandomKeys()
	if err != nil {
		return nil, err
	}

	user := &regressionTestUser{
		XPriv:   keys.XPriv(),
		XPub:    keys.XPub(),
		Paymail: preparePaymail(paymail, paymailDomain),
	}

	cfg := config.NewDefaultConfig(instanceUrl)
	adminClient, err := wallet.NewAdminAPIWithXPriv(cfg, adminXPriv)
	if err != nil {
		return nil, err
	}

	_, err = adminClient.CreateXPub(ctx, &commands.CreateUserXpub{
		Metadata: map[string]any{"some_metadata": "remove"},
		XPub:     user.XPub,
	})
	if err != nil {
		return nil, err
	}

	_, err = adminClient.CreatePaymail(ctx, &commands.CreatePaymail{
		Metadata:   map[string]any{},
		Key:        user.XPub,
		Address:    user.Paymail,
		PublicName: "Regression tests",
		Avatar:     "",
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}

// removeRegisteredPaymail soft deletes paymail from the SPV Wallet.
func removeRegisteredPaymail(ctx context.Context, paymail string, instanceUrl string, adminXPriv string) error {
	cfg := config.NewDefaultConfig(instanceUrl)
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
	cfg := config.NewDefaultConfig(fromInstance)
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
	cfg := config.NewDefaultConfig(fromInstance)
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
	cfg := config.NewDefaultConfig(fromInstance)
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
