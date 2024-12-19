//go:build regression
// +build regression

package regressiontests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
)

func TestRegressionWorkflow(t *testing.T) {
	ctx := context.Background()
	spvWalletPG, spvWalletSL := initSPVWalletAPIs(t)

	t.Log("Step 1: Setup success: created SPV client instances with test users")
	t.Logf("SPV clients for env: %s, user: %s, admin: %s, leader: %s", spvWalletPG.cfg.envURL, spvWalletPG.user.alias, spvWalletPG.admin.alias, spvWalletPG.leader.alias)
	t.Logf("SPV clients for env: %s, user: %s, admin: %s, leader: %s", spvWalletSL.cfg.envURL, spvWalletSL.user.alias, spvWalletSL.admin.alias, spvWalletSL.leader.alias)

	t.Run("Step 2: The leader user attempts to retrieve the shared configuration response from their SPV Wallet API instance.", func(t *testing.T) {
		tests := []struct {
			name                string
			leader              *user
			user                *user
			admin               *admin
			expectedPaymailsLen int
		}{
			{
				name:   fmt.Sprintf("%s should set paymail domain after fetching shared config", spvWalletPG.leader.alias),
				leader: spvWalletPG.leader,
				user:   spvWalletPG.user,
				admin:  spvWalletPG.admin,
			},
			{
				name:   fmt.Sprintf("%s should set paymail domain after fetching shared config", spvWalletSL.leader.alias),
				leader: spvWalletSL.leader,
				user:   spvWalletSL.user,
				admin:  spvWalletSL.admin,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				// when:
				got, err := tc.leader.client.SharedConfig(ctx)

				// then:
				if err != nil {
					t.Errorf("Shared config wasn't successful retrieve by %s. Expect to get nil error, got error: %v", tc.leader.paymail, err)
				}

				if len(got.PaymailDomains) != 1 {
					t.Errorf("Shared config retrieved by %s should have single paymail domain. Got: %d paymail domains", tc.leader.paymail, len(got.PaymailDomains))
				}
				domain := got.PaymailDomains[0]
				if len(domain) == 0 {
					t.Errorf("Shared config retrieved by %s should not be empty string", tc.leader.paymail)
				}

				tc.leader.setPaymail(domain)
				tc.admin.setPaymail(domain)
				tc.user.setPaymail(domain)
			})
		}
	})

	t.Run("Step 3: The SPV Wallet admin client attempt to add xPub records of the user from the same env by making request to their SPV Wallet API instance", func(t *testing.T) {
		tests := []struct {
			name  string
			user  *user
			admin *admin
		}{
			{
				name:  fmt.Sprintf("%s should add xPub record %s for %s", spvWalletPG.admin.paymail, spvWalletPG.user.xPub, spvWalletPG.user.paymail),
				admin: spvWalletPG.admin,
				user:  spvWalletPG.user,
			},
			{
				name:  fmt.Sprintf("%s should add xPub record %s for %s", spvWalletPG.admin.paymail, spvWalletSL.user.xPub, spvWalletSL.user.paymail),
				admin: spvWalletSL.admin,
				user:  spvWalletSL.user,
			},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				// when:
				xPub, err := tc.admin.client.CreateXPub(ctx, &commands.CreateUserXpub{XPub: tc.user.xPub})
				// then:
				if err != nil {
					t.Errorf("xPub record %s wasn't created successfully for %s by %s. Got error: %v", tc.user.xPub, tc.user.paymail, tc.admin.paymail, err)
				}

				if xPub == nil {
					t.Errorf("Expected to get non-nil xPub response after sending creation request by %s.", tc.admin.paymail)
				}

				if xPub != nil && err == nil {
					t.Logf("xPub record %s was created successfully for %s by %s", tc.user.xPub, tc.user.paymail, tc.admin.paymail)
				}
			})
		}
	})

	t.Run("Step 4: The SPV Wallet admin clients attempt to add paymail record of the user from the same env by making request to their SPV Wallet API instance", func(t *testing.T) {
		tests := []struct {
			name  string
			user  *user
			admin *admin
		}{
			{
				name:  fmt.Sprintf("%s should add paymail record %s for the user %s", spvWalletPG.admin.paymail, spvWalletPG.user.paymail, spvWalletPG.user.alias),
				admin: spvWalletPG.admin,
				user:  spvWalletPG.user,
			},
			{
				name:  fmt.Sprintf("%s should add paymail record %s for the user %s", spvWalletPG.admin.paymail, spvWalletSL.user.paymail, spvWalletSL.user.alias),
				admin: spvWalletSL.admin,
				user:  spvWalletSL.user,
			},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				// when:
				paymail, err := tc.admin.client.CreatePaymail(ctx, &commands.CreatePaymail{
					Key:        tc.user.xPub,
					Address:    tc.user.paymail,
					PublicName: "Regression tests",
				})

				// then:
				if err != nil {
					t.Errorf("Paymail record %s wasn't created successfully for %s by %s. Got error: %v", tc.user.paymail, tc.user.alias, tc.admin.paymail, err)
				}

				if paymail == nil {
					t.Errorf("Expected to get non-nil paymail addresss response after sending creation request by %s.", tc.admin.paymail)
				}

				if err == nil && paymail != nil {
					t.Logf("Paymail record %s was created successfully for %s by %s.", tc.user.paymail, tc.user.alias, tc.admin.paymail)
				}
			})
		}
	})

	t.Run("Step 5: The leaders from one env attempts to transfer funds to users from another env using the appropriate SPV Wallet API instance", func(t *testing.T) {
		tests := []struct {
			name   string
			leader *user
			user   *user
			funds  uint64
		}{
			{
				leader: spvWalletPG.leader,
				user:   spvWalletSL.user,
				funds:  3,
				name:   fmt.Sprintf("%s should transfer 3 satoshis to the user %s", spvWalletPG.leader.paymail, spvWalletSL.user.paymail),
			},
			{
				leader: spvWalletSL.leader,
				user:   spvWalletPG.user,
				funds:  2,
				name:   fmt.Sprintf("%s should transfer 2 satoshis to the user %s", spvWalletSL.leader.paymail, spvWalletPG.user.paymail),
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				// given:
				recipientBalance := checkBalance(ctx, t, tc.user.client)
				leaderBalance := checkBalance(ctx, t, tc.leader.client)
				if leaderBalance < tc.funds {
					t.Fatalf("Transfer funds %d wasn't successful from  %s to %s. Due to insufficient balance. Need to have at least: %d sathoshis. Got: %d",
						tc.funds, tc.leader.paymail, tc.user.paymail, tc.funds, leaderBalance)
				}

				// when:
				transaction, err := tc.leader.transferFunds(ctx, tc.user.paymail, tc.funds)

				// then:
				if err != nil {
					t.Errorf("Transfer funds %d wasn't successful from %s to %s. Expect to get nil error after making transaction, got error: %v", tc.funds, tc.leader.paymail, tc.user.paymail, err)
				}

				if transaction == nil {
					t.Errorf("Expected to get non-nil transaction response after transfer funds %d from %s to %s", tc.funds, tc.leader.paymail, tc.user.paymail)
				}

				if transaction != nil && err == nil {
					transferBalance := transferBalance{
						sender:           tc.leader,
						recipient:        tc.user,
						senderBalance:    leaderBalance,
						recipientBalance: recipientBalance,
						transactionID:    transaction.ID,
						fee:              transaction.Fee,
						funds:            tc.funds,
					}

					transferBalance.check(ctx, t)
				}
			})
		}
	})

	t.Run("Step 6: The user from one env attempts to transfer funds to the user from external env using the appropriate SPV Wallet API instance", func(t *testing.T) {
		tests := []struct {
			name      string
			sender    *user
			recipient *user
			funds     uint64
		}{
			{
				name:      fmt.Sprintf("%s should transfer 2 satoshis to %s", spvWalletSL.user.paymail, spvWalletPG.user.paymail),
				sender:    spvWalletSL.user,
				recipient: spvWalletPG.user,
				funds:     2,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				// given:
				recipientBalance := checkBalance(ctx, t, tc.recipient.client)
				senderBalance := checkBalance(ctx, t, tc.sender.client)
				if senderBalance < tc.funds {
					t.Fatalf("Transfer funds %d wasn't successful from sender %s to recipient %s. Due to insufficient balance. Need to have at least: %d sathoshis. Got: %d",
						tc.funds, tc.sender.paymail, tc.recipient.paymail, tc.funds, senderBalance)
				}

				// when:
				transaction, err := tc.sender.transferFunds(ctx, tc.recipient.paymail, tc.funds)

				// then:
				if err != nil {
					t.Errorf("Transfer funds %d wasn't successful from sender %s to recipient %s. Expect to get nil error after making transaction, got error: %v",
						tc.funds, tc.sender.paymail, tc.recipient.paymail, err)
				}

				if transaction == nil {
					t.Errorf("Expected to get non-nil transaction response after transfer funds %d from sender %s to recipient %s", tc.funds, tc.sender.paymail, tc.recipient.paymail)
				}

				if err == nil && transaction != nil {
					transferBalance := transferBalance{
						sender:           tc.sender,
						recipient:        tc.recipient,
						senderBalance:    senderBalance,
						recipientBalance: recipientBalance,
						transactionID:    transaction.ID,
						fee:              transaction.Fee,
						funds:            tc.funds,
					}

					transferBalance.check(ctx, t)
				}
			})
		}
	})

	t.Run("Step 7: The leader attempts to remove created actor paymails using the appropriate SPV Wallet API instance.", func(t *testing.T) {
		tests := []struct {
			name    string
			admin   *admin
			paymail string
		}{
			{
				name:    fmt.Sprintf("%s should delete %s paymail record", spvWalletPG.admin.paymail, spvWalletPG.user.paymail),
				admin:   spvWalletPG.admin,
				paymail: spvWalletPG.user.paymail,
			},
			{
				name:    fmt.Sprintf("%s should delete %s paymail record", spvWalletSL.admin.paymail, spvWalletSL.user.paymail),
				admin:   spvWalletSL.admin,
				paymail: spvWalletSL.user.paymail,
			},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				// when:
				err := tc.admin.client.DeletePaymail(ctx, tc.paymail)

				// then:
				if err != nil {
					t.Errorf("Delete paymail %s wasn't successful by %s. Expect to get nil error, got error: %v", tc.paymail, tc.admin.paymail, err)
				}

				if err == nil {
					t.Logf("Delete paymail %s was successful by %s", tc.paymail, tc.admin.paymail)
				}
			})
		}
	})
}

func initSPVWalletAPIs(t *testing.T) (*spvWalletAPI, *spvWalletAPI) {
	const adminXPriv = "xprv9s21ZrQH143K3CbJXirfrtpLvhT3Vgusdo8coBritQ3rcS7Jy7sxWhatuxG5h2y1Cqj8FKmPp69536gmjYRpfga2MJdsGyBsnB12E19CESK"

	const (
		clientOneURL         = "CLIENT_ONE_URL"
		clientOneLeaderXPriv = "CLIENT_ONE_LEADER_XPRIV"
		clientTwoURL         = "CLIENT_TWO_URL"
		clientTwoLeaderXPriv = "CLIENT_TWO_LEADER_XPRIV"
	)

	const (
		alias1 = "UserSLRegressionTest"
		alias2 = "UserPGRegressionTest"
	)

	spvWalletSL, err := initSPVWalletAPI(&spvWalletAPIConfig{
		envURL:     lookupEnvOrDefault(t, clientOneURL, ""),
		envXPriv:   lookupEnvOrDefault(t, clientOneLeaderXPriv, ""),
		adminXPriv: adminXPriv,
	}, alias1)
	if err != nil {
		t.Fatalf("Step 1: Setup failed could not initialize the clients for env: %s. Got error: %v", spvWalletSL.cfg.envURL, err)
	}

	spvWalletPG, err := initSPVWalletAPI(&spvWalletAPIConfig{
		envURL:     lookupEnvOrDefault(t, clientTwoURL, ""),
		envXPriv:   lookupEnvOrDefault(t, clientTwoLeaderXPriv, ""),
		adminXPriv: adminXPriv,
	}, alias2)
	if err != nil {
		t.Fatalf("Step 1: Setup failed could not initialize the clients for env: %s. Got error: %v", spvWalletPG.cfg.envURL, err)
	}

	return spvWalletPG, spvWalletSL
}
