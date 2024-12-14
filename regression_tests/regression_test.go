//go:build regression
// +build regression

package regressiontests

import (
	"context"
	"os"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/regression_tests/services"
	"github.com/stretchr/testify/require"
)

func TestRegression_TransactionsWorkflow(t *testing.T) {
	const adminXPriv = "xprv9s21ZrQH143K3CbJXirfrtpLvhT3Vgusdo8coBritQ3rcS7Jy7sxWhatuxG5h2y1Cqj8FKmPp69536gmjYRpfga2MJdsGyBsnB12E19CESK"
	const actor1Alias = "Actor1RegressionTest"
	const actor2Alias = "Actor2RegressionTest"

	ctx := context.Background()
	leaderService1, actorService1 := initTestServices(ctx, t, &services.LeaderServiceConfig{
		EnvURL:     lookupEnvOrDefault("CLIENT_ONE_URL", ""),
		EnvXPriv:   lookupEnvOrDefault("CLIENT_ONE_LEADER_XPRIV", ""),
		AdminXPriv: adminXPriv,
	}, actor1Alias)

	leaderService2, actorService2 := initTestServices(ctx, t, &services.LeaderServiceConfig{
		EnvURL:     lookupEnvOrDefault("CLIENT_TWO_URL", ""),
		EnvXPriv:   lookupEnvOrDefault("CLIENT_TWO_LEADER_XPRIV", ""),
		AdminXPriv: adminXPriv,
	}, actor2Alias)

	defer t.Cleanup(func() {
		err := leaderService1.RemoveActors(ctx)
		if err != nil {
			t.Errorf("Cleanup failed: could not actors by  leader service %s: %v", leaderService1.Name(), err)
		}
	})

	defer t.Cleanup(func() {
		err := leaderService2.RemoveActors(ctx)
		if err != nil {
			t.Errorf("Cleanup failed: could not actors by leader service %s: %v", leaderService1.Name(), err)
		}
	})

	actor1 := actorService1.Actor()
	actor2 := actorService2.Actor()

	const minimalTransactionOutputValue = int64(-1)

	// First transfer from leader2 to actor1
	const firstTransferFunds = 3

	leader2Transaction, err := leaderService2.TransferFunds(ctx, actor1.Paymail(), firstTransferFunds)
	require.NoError(t, err)
	require.GreaterOrEqual(t, minimalTransactionOutputValue, leader2Transaction.OutputValue)

	expectedActor1Balance := uint64(1)
	actor1Balance, err := actorService1.Balance(ctx)
	require.NoError(t, err)
	require.GreaterOrEqual(t, actor1Balance, expectedActor1Balance)

	expectedActor1TransactionsCount := 1
	actor1Transactions, err := actorService1.Transactions(ctx)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(actor1Transactions), expectedActor1TransactionsCount)

	// Second transfer from leader1 to actor2
	const secondTransferFunds = 2

	leader1Transaction, err := leaderService1.TransferFunds(ctx, actor2.Paymail(), secondTransferFunds)
	require.NoError(t, err)
	require.GreaterOrEqual(t, minimalTransactionOutputValue, leader1Transaction.OutputValue)

	expectedActor2Balance := uint64(1)
	actor2Balance, err := actorService2.Balance(ctx)
	require.NoError(t, err)
	require.GreaterOrEqual(t, actor2Balance, expectedActor2Balance)

	expectedActor2TransactionsCount := 1
	actor2Transactions, err := actorService2.Transactions(ctx)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(actor2Transactions), expectedActor2TransactionsCount)

	// Third transfer from actor1 to actor2
	const thirdTransferFunds = 2

	actor1Transaction, err := actorService1.TransferFunds(ctx, actor2.Paymail(), thirdTransferFunds)
	require.NoError(t, err)
	require.GreaterOrEqual(t, minimalTransactionOutputValue, actor1Transaction.OutputValue)

	minimalActor1Balance := uint64(0)
	actor1Balance, err = actorService1.Balance(ctx)
	require.NoError(t, err)
	require.GreaterOrEqual(t, actor1Balance, minimalActor1Balance)

	expectedActor1TransactionsCount = 2
	actor1Transactions, err = actorService1.Transactions(ctx)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(actor1Transactions), expectedActor1TransactionsCount)

	expectedActor2Balance = uint64(2)
	actor2Balance, err = actorService2.Balance(ctx)
	require.NoError(t, err)
	require.GreaterOrEqual(t, actor2Balance, expectedActor2Balance)

	expectedActor2TransactionsCount = 2
	actor2Transactions, err = actorService2.Transactions(ctx)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(actor2Transactions), expectedActor2TransactionsCount)
}

func initTestServices(ctx context.Context, t *testing.T, cfg *services.LeaderServiceConfig, alias string) (*services.LeaderService, *services.ActorService) {
	t.Helper()

	leaderService, err := services.NewLeaderService(cfg)
	if err != nil {
		t.Fatalf("Setup failed: could not initialize leader service: %v", err)
	}
	actor, err := leaderService.CreateActor(ctx, alias)
	if err != nil {
		t.Fatalf("Setup failed: could not create actor: %v", err)
	}
	actorService, err := services.NewActorService(cfg.EnvURL, actor)
	if err != nil {
		t.Fatalf("Setup failed: count not initialize actor service: %v", err)
	}
	return leaderService, actorService
}

func lookupEnvOrDefault(env string, s string) string {
	v, ok := os.LookupEnv(env)
	if ok {
		return v
	}
	return s
}
