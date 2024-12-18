package regressiontests

import (
	"context"
	"fmt"
	"os"
	"testing"
)

func TestRegressionWorkflow(t *testing.T) {
	ctx := context.Background()
	spvWalletPG, spvWalletSL := setupTest(t)

	t.Run("Step 1: The SPV Wallet admin clients attempt to add xPub records by making requests to their SPV Wallet API instance.", func(t *testing.T) {
		tests := []struct {
			name string
			xPub string
			spv  *spvWallet
		}{
			{
				name: fmt.Sprintf("%s admin client should add xPub record %s for the user %s.", spvWalletPG.name, spvWalletPG.user.xPub, spvWalletPG.user.alias),
				spv:  spvWalletPG,
				xPub: spvWalletPG.user.xPub,
			},
			{
				name: fmt.Sprintf("%s admin client should add xPub record %s for the user %s.", spvWalletPG.name, spvWalletSL.user.xPub, spvWalletSL.user.alias),
				spv:  spvWalletSL,
				xPub: spvWalletSL.user.xPub,
			},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				// when:
				err := tc.spv.admin.createXPub(ctx, tc.xPub)

				// then:
				if err != nil {
					t.Errorf("xPub record: %s wasn't created successfully for the user %s with paymail %s by %s. Got error: %v", tc.spv.user.xPub, tc.spv.user.alias, tc.spv.user.paymail, tc.name, err)
				} else {
					t.Logf("xPub record: %s was created successfully for the user %s with paymail %s by %s.", tc.spv.user.xPub, tc.spv.user.alias, tc.spv.user.paymail, tc.name)
				}
			})
		}
	})

	t.Run("Step 2: The SPV Wallet admin clients attempt to add paymail records by making requests to their SPV Wallet API instance.", func(t *testing.T) {
		tests := []struct {
			name    string
			paymail string
			xPub    string
			spv     *spvWallet
		}{
			{
				name:    fmt.Sprintf("%s admin client should add paymail record %s for user %s", spvWalletPG.name, spvWalletPG.user.paymail, spvWalletPG.user.alias),
				spv:     spvWalletPG,
				paymail: spvWalletPG.user.paymail,
				xPub:    spvWalletPG.user.xPub,
			},
			{
				name:    fmt.Sprintf("%s admin client should add paymail record %s for user %s", spvWalletPG.name, spvWalletSL.user.paymail, spvWalletSL.user.alias),
				spv:     spvWalletSL,
				paymail: spvWalletSL.user.paymail,
				xPub:    spvWalletSL.user.xPub,
			},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				// when:
				err := tc.spv.admin.createPaymail(ctx, tc.xPub, tc.paymail)

				// then:
				if err != nil {
					t.Errorf("Paymail record: %s wasn't created successfully for user %s with paymail %s by %s. Got error: %v", tc.spv.user.paymail, tc.spv.user.alias, tc.spv.user.paymail, tc.name, err)
				} else {
					t.Logf("Paymail record: %s was created successfully for the user %s with paymail %s by %s.", tc.spv.user.paymail, tc.spv.user.alias, tc.spv.user.paymail, tc.name)
				}
			})
		}
	})
}

func setupTest(t *testing.T) (*spvWallet, *spvWallet) {
	t.Helper()

	const admin = "xprv9s21ZrQH143K3CbJXirfrtpLvhT3Vgusdo8coBritQ3rcS7Jy7sxWhatuxG5h2y1Cqj8FKmPp69536gmjYRpfga2MJdsGyBsnB12E19CESK"

	const (
		clientOneURL         = "CLIENT_ONE_URL"
		clientOneLeaderXPriv = "CLIENT_ONE_LEADER_XPRIV"
		clientTwoURL         = "CLIENT_TWO_URL"
		clientTwoLeaderXPriv = "CLIENT_TWO_LEADER_XPRIV"
	)

	const (
		alias1 = "Actor1RegressionTest"
		alias2 = "Actor2RegressionTest"
	)

	client1, err := initSPVWallet(&spvWalletConfig{
		envURL:     lookupEnvOrDefault(clientOneURL, ""),
		envXPriv:   lookupEnvOrDefault(clientOneLeaderXPriv, ""),
		adminXPriv: admin,
		name:       "SPV Wallet PG",
	}, alias1)
	if err != nil {
		t.Fatalf("Setup failed: could not initialize SPV Wallet PG: %v", err)
	}

	client2, err := initSPVWallet(&spvWalletConfig{
		envURL:     lookupEnvOrDefault(clientTwoURL, ""),
		envXPriv:   lookupEnvOrDefault(clientTwoLeaderXPriv, ""),
		adminXPriv: admin,
		name:       "SPV Wallet SL",
	}, alias2)
	if err != nil {
		t.Fatalf("Setup failed: could not initialize SPV Wallet SL: %v", err)
	}

	return client1, client2
}

func lookupEnvOrDefault(env string, s string) string {
	v, ok := os.LookupEnv(env)
	if ok {
		return v
	}
	return s
}
