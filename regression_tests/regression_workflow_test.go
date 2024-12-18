//go:build regression
// +build regression

package regressiontests

import (
	"context"
	"fmt"
	"os"
	"testing"

	wallet "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/stretchr/testify/require"
)

func TestRegressionWorkflow(t *testing.T) {
	const (
		clientOneURL         = "CLIENT_ONE_URL"
		clientOneLeaderXPriv = "CLIENT_ONE_LEADER_XPRIV"
		clientTwoURL         = "CLIENT_TWO_URL"
		clientTwoLeaderXPriv = "CLIENT_TWO_LEADER_XPRIV"
	)

	const (
		admin = "xprv9s21ZrQH143K3CbJXirfrtpLvhT3Vgusdo8coBritQ3rcS7Jy7sxWhatuxG5h2y1Cqj8FKmPp69536gmjYRpfga2MJdsGyBsnB12E19CESK"
	)

	ctx := context.Background()
	leader1 := initLeader(t, &leaderConfig{
		envURL:     lookupEnvOrDefault(clientOneURL, ""),
		envXPriv:   lookupEnvOrDefault(clientOneLeaderXPriv, ""),
		adminXPriv: admin,
	})
	leader2 := initLeader(t, &leaderConfig{
		envURL:     lookupEnvOrDefault(clientTwoURL, ""),
		envXPriv:   lookupEnvOrDefault(clientTwoLeaderXPriv, ""),
		adminXPriv: admin,
	})

	t.Run("Step 1: The leader attempts to retrieve the shared configuration response from their SPV Wallet API instance.", func(t *testing.T) {
		tests := []struct {
			name                string
			alias               string
			leader              *leader
			expectedPaymailsLen int
		}{
			{
				name:                fmt.Sprintf("Leader [%s] should set paymail domain after fetching shared config", leader1.name()),
				leader:              leader1,
				alias:               "Actor1RegressionTest",
				expectedPaymailsLen: 1,
			},
			{
				name:                fmt.Sprintf("Leader [%s] should set paymail domain after fetching shared config", leader2.name()),
				leader:              leader2,
				alias:               "Actor2RegressionTest",
				expectedPaymailsLen: 1,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				// given:
				got, err := tc.leader.userAPI.SharedConfig(ctx)

				// when:
				if err != nil {
					t.Errorf("[Leader: %s] expected to get nil err after UserAPI SharedConfig call; got: %v", tc.leader.name(), err)
				}
				if len(got.PaymailDomains) != 1 {
					t.Errorf("[Leader: %s] expected to have single paymail domain in the paymail domains slice; got: %d paymail domain", tc.leader.name(), len(got.PaymailDomains))
				}
				domain := got.PaymailDomains[0]
				if len(domain) == 0 {
					t.Errorf("[Leader: %s] expected have non empty string as paymail domain", tc.leader.name())
				}

				tc.leader.setPaymailDomain(domain)
				tc.leader.setActor(initActor(t, tc.alias, domain))
			})
		}
	})

	t.Run("Step 2: The leader attempts to add xpub records by making requests to their SPV Wallet API instance.", func(t *testing.T) {
		tests := []struct {
			name   string
			leader *leader
		}{
			{
				name:   fmt.Sprintf("Leader [%s] should add xpub record for Actor [%s].", leader1.name(), leader1.actor.paymail),
				leader: leader1,
			},
			{
				name:   fmt.Sprintf("Leader [%s] should add xpub record for Actor [%s].", leader2.name(), leader2.actor.paymail),
				leader: leader2,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				// given:
				actor := tc.leader.actor
				got, err := tc.leader.adminAPI.CreateXPub(ctx, actor.createUserXpubCommand())

				// when:
				if err != nil {
					t.Errorf("[Leader: %s][Actor: %s] expected to get nil err after CreateXPub AdminAPI call; got: %v", tc.leader.name(), actor.paymail, err)
				}
				if got != nil {
					t.Errorf("[Leader: %s][Actor: %s] expected to get non nil Xpub response after CreateXPub AdminAPI call", tc.leader.name(), actor.paymail)
				}
			})
		}
	})

	t.Run("Step 3: The leader attempts to add paymail records by making requests to their SPV Wallet API instance.", func(t *testing.T) {
		tests := []struct {
			name   string
			leader *leader
		}{
			{
				name:   fmt.Sprintf("Leader [%s] should add paymail domain record for Actor[%s]", leader1.name(), leader1.actor.paymail),
				leader: leader1,
			},
			{
				name:   fmt.Sprintf("Leader [%s] should add paymail domain record for Actor[%s]", leader2.name(), leader2.actor.paymail),
				leader: leader2,
			},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				// given:
				actor := tc.leader.actor
				got, err := tc.leader.adminAPI.CreatePaymail(ctx, actor.createPaymailCommand())

				// when:
				if err != nil {
					t.Errorf("[Leader: %s][Actor: %s] expected to get nil err after CreatePaymail AdminAPI call; got: %v", tc.leader.name(), actor.paymail, err)
				}
				if got != nil {
					t.Errorf("[Leader: %s][Actor: %s] expected to get non nil Xpub response after CreatePaymail AdminAPI call", tc.leader.name(), actor.paymail)
				}
			})
		}
	})

	t.Run("Step 4: The leader attempts to transfer funds to external actors using the appropriate SPV Wallet API instance.", func(t *testing.T) {
		tests := []struct {
			name      string
			sender    *leader
			recipient *leader
			funds     uint64
		}{
			{
				name:      fmt.Sprintf("Leader [%s] should transfer 3 satoshis to Actor [%s]", leader2.name(), leader1.actor.paymail),
				sender:    leader2,
				recipient: leader1,
				funds:     3,
			},
			{
				name:      fmt.Sprintf("Leader [%s] should transfer 2 satoshis to Actor [%s]", leader1.name(), leader2.actor.paymail),
				sender:    leader1,
				recipient: leader2,
				funds:     2,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				// given:
				prevRecipientBalance := checkBalance(ctx, t, tc.recipient.userAPI)
				prevSenderBalance := checkBalance(ctx, t, tc.sender.userAPI)
				require.GreaterOrEqual(t, prevSenderBalance, tc.funds)

				// when:
				to := tc.recipient.actor.paymail
				transaction, err := tc.sender.userAPI.SendToRecipients(ctx, &commands.SendToRecipients{
					Recipients: []*commands.Recipients{
						{
							To:       to,
							Satoshis: tc.funds,
						},
					},
					Metadata: map[string]any{"description": "regression-test"},
				})

				// then:
				if err != nil {
					t.Errorf("[Sender: %s][Recipient: %s] expected to get nil err after SendToRecipients UserAPI call; got: %v", tc.sender.name(), to, err)
				}

				currentSenderBalance := checkBalance(ctx, t, tc.sender.userAPI)
				expectedSenderBalance := prevSenderBalance - tc.funds - transaction.Fee
				if currentSenderBalance != expectedSenderBalance {
					t.Errorf("[Sender: %s] expected to get reduced balance from: %d to: %d after making the transaction; got: %d", tc.sender.name(), prevRecipientBalance, expectedSenderBalance, currentSenderBalance)
				}

				currentRecipientBalance := checkBalance(ctx, t, tc.recipient.userAPI)
				expectedRecipientBalance := prevRecipientBalance + tc.funds
				if currentRecipientBalance != expectedRecipientBalance {
					t.Errorf("[Recipient: %s] expected to get increased balance from: %d to: %d after making the transaction; got: %d", tc.recipient.name(), prevRecipientBalance, expectedRecipientBalance, currentRecipientBalance)
				}
			})
		}
	})

	t.Run("Step 5: The leader attempts to remove created actor paymails using the appropriate SPV Wallet API instance.", func(t *testing.T) {
		tests := []struct {
			name   string
			leader *leader
		}{
			{
				name:   fmt.Sprintf("Leader [%s] should delete Actor [%s] paymail record", leader2.name(), leader1.actor.paymail),
				leader: leader1,
			},
			{
				name:   fmt.Sprintf("Leader [%s] should delete Actor [%s] paymail record", leader2.name(), leader2.actor.paymail),
				leader: leader2,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				// when:
				paymail := tc.leader.actor.paymail
				err := tc.leader.adminAPI.DeletePaymail(ctx, tc.leader.actor.paymail)
				if err != nil {
					t.Errorf("[Leader: %s][Paymail: %s] expected to get nil err after DeletePaymail AdminAPI call; got: %v", tc.leader.name(), paymail, err)
				}
			})
		}
	})
}

func checkBalance(ctx context.Context, t *testing.T, api *wallet.UserAPI) uint64 {
	t.Helper()
	xPub, err := api.XPub(ctx)
	if err != nil {
		t.Fatalf("Balance check failed: could not fetch the balance %v", err)
	}
	if xPub == nil {
		t.Fatal("Balance check failed: expected to get non nil response after UserAPI XPub call")
	}
	return xPub.CurrentBalance
}

func initLeader(t *testing.T, cfg *leaderConfig) *leader {
	t.Helper()
	leader, err := NewLeader(cfg)
	if err != nil {
		t.Fatalf("Setup failed: could not initialize the leader[env: %s]: %v", cfg.envURL, err)
	}
	return leader
}

func initActor(t *testing.T, alias, domain string) *actor {
	t.Helper()
	actor, err := newActor(alias, domain)
	if err != nil {
		t.Fatalf("Setup failed: could not initialize the actor [alias: %s, domain: %s]: %v", alias, domain, err)
	}
	return actor
}

func lookupEnvOrDefault(env string, s string) string {
	v, ok := os.LookupEnv(env)
	if ok {
		return v
	}
	return s
}
