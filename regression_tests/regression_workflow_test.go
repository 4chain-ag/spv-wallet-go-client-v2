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
		alias1 = "Actor1RegressionTest"
		alias2 = "Actor2RegressionTest"
		admin  = "xprv9s21ZrQH143K3CbJXirfrtpLvhT3Vgusdo8coBritQ3rcS7Jy7sxWhatuxG5h2y1Cqj8FKmPp69536gmjYRpfga2MJdsGyBsnB12E19CESK"
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
			leader              *leader
			expectedErr         error
			expectedPaymailsLen int
		}{
			{
				name:                fmt.Sprintf("Leader [%s] should set paymail domain after fetching shared config", leader1.name()),
				leader:              leader1,
				expectedPaymailsLen: 1,
			},
			{
				name:                fmt.Sprintf("Leader [%s] should set paymail domain after fetching shared config", leader2.name()),
				leader:              leader2,
				expectedPaymailsLen: 1,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				// given:
				got, err := tc.leader.userAPI.SharedConfig(ctx)

				// when:
				require.ErrorIs(t, err, tc.expectedErr)
				require.Len(t, got.PaymailDomains, tc.expectedPaymailsLen)
				require.NotEmpty(t, got, got.PaymailDomains[0])

				tc.leader.setPaymailDomain(got.PaymailDomains[0])
			})
		}
	})

	leader1.setActor(initActor(t, alias1, leader1.domain))
	leader2.setActor(initActor(t, alias2, leader2.domain))

	t.Run("Step 2: The leader attempts to add paymail records by making requests to their SPV Wallet API instance.", func(t *testing.T) {
		tests := []struct {
			name        string
			leader      *leader
			expectedErr error
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
				require.ErrorIs(t, err, tc.expectedErr)
				require.NotNil(t, got)
			})
		}
	})
	t.Run("Step 3: The leader attempts to add xpub records by making requests to their SPV Wallet API instance.", func(t *testing.T) {
		tests := []struct {
			name        string
			leader      *leader
			expectedErr error
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
				require.ErrorIs(t, err, tc.expectedErr)
				require.NotNil(t, got)
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
				require.NoError(t, err)

				currentSenderBalance := checkBalance(ctx, t, tc.sender.userAPI)
				currentRecipientBalance := checkBalance(ctx, t, tc.recipient.userAPI)

				expectedSenderBalance := prevSenderBalance - tc.funds - transaction.Fee
				expectedRecipientBalance := prevRecipientBalance + tc.funds

				require.Equal(t, expectedSenderBalance, currentSenderBalance)
				require.Equal(t, expectedRecipientBalance, currentRecipientBalance)
			})
		}
	})

	t.Run("Step 5: Actor1 from the Leader1 instance attempts to transfer funds to Actor2 from the Leader2 instance.", func(t *testing.T) {
		tests := []struct {
			name      string
			sender    *leader
			recipient *leader
			funds     uint64
		}{
			{
				name:      fmt.Sprintf("Actor [%s] should transfer 2 satoshis to actor: %s.", leader1.actor.paymail, leader1.actor.paymail),
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
				require.NoError(t, err)

				currentSenderBalance := checkBalance(ctx, t, tc.sender.userAPI)
				currentRecipientBalance := checkBalance(ctx, t, tc.recipient.userAPI)

				expectedSenderBalance := prevSenderBalance - tc.funds - transaction.Fee
				expectedRecipientBalance := prevRecipientBalance + tc.funds

				require.Equal(t, expectedSenderBalance, currentSenderBalance)
				require.Equal(t, expectedRecipientBalance, currentRecipientBalance)
			})
		}
	})

	t.Run("Step 6: The leader attempts to remove created actor paymails using the appropriate SPV Wallet API instance.", func(t *testing.T) {
		tests := []struct {
			name        string
			leader      *leader
			expectedErr error
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
				err := tc.leader.adminAPI.DeletePaymail(ctx, tc.leader.actor.paymail)
				require.NoError(t, err)
			})
		}
	})
}

func checkBalance(ctx context.Context, t *testing.T, api *wallet.UserAPI) uint64 {
	t.Helper()
	xPub, err := api.XPub(ctx)
	require.NoError(t, err)
	require.NotNil(t, xPub)
	return xPub.CurrentBalance
}

func initLeader(t *testing.T, cfg *leaderConfig) *leader {
	t.Helper()
	leader, err := NewLeader(cfg)
	require.NoError(t, err)
	return leader
}

func initActor(t *testing.T, alias, domain string) *actor {
	t.Helper()
	actor, err := newActor(alias, domain)
	require.NoError(t, err)
	return actor
}

func lookupEnvOrDefault(env string, s string) string {
	v, ok := os.LookupEnv(env)
	if ok {
		return v
	}
	return s
}
