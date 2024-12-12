package regressiontests

import (
	"context"
	"testing"
)

func TestRegressionWorkflow(t *testing.T) {
	const (
		actor1Alias = "actor1111111111_regression_tests"
		actor2Alias = "actor2222222222_regression_tests" // TODO: To rename.
	)

	ctx := context.Background()
	leaderService1, actorService1 := initTestServices(ctx, "CLIENT_ONE_URL", "CLIENT_ONE_LEADER_XPRIV", actor1Alias, t)
	leaderService2, actorService2 := initTestServices(ctx, "CLIENT_TWO_URL", "CLIENT_TWO_LEADER_XPRIV", actor2Alias, t)

	defer t.Cleanup(func() {
		err := leaderService1.RemoveActors(ctx)
		if err != nil {
			t.Fatalf("Cleanup failed: could not actors in leader service %s: %v", leaderService1.name, err)
		}
	})

	defer t.Cleanup(func() {
		err := leaderService2.RemoveActors(ctx)
		if err != nil {
			t.Fatalf("Cleanup failed: could not actors in leader service %s: %v", leaderService2.name, err)
		}
	})

	TransferFundsFromLeaderToActor(ctx, t, leaderService2, actorService1, 3)
	TransferFundsFromLeaderToActor(ctx, t, leaderService1, actorService2, 2)

	TransferFundsFromActorToPaymail(ctx, t, actorService1, actorService2.actor.paymail)

	CheckActorStateAfterFundsTransfer(ctx, t, actorService2, 2, 2)
}

func CheckActorStateAfterFundsTransfer(ctx context.Context, t *testing.T, actorService *ActorService, minimumBalance uint64, minimumTransactions int) {
	balance, err := actorService.Balance(ctx)
	if err != nil {
		t.Fatalf("ActorService - Balance failed: could not check balance status: %v", err)
	}
	if balance < minimumBalance {
		t.Fatalf("expected to get balance value greater or equal %d, got: %d", minimumBalance, balance)
	}

	transactions, err := actorService.Transactions(ctx)
	if err != nil {
		t.Fatalf("ActorService - Transactions failed: coult not get transactions: %v", err)
	}
	if transactions == nil {
		t.Fatalf("expected to get non nil transactions slice for %s paymail", actorService.actor.paymail)
	}
	if len(transactions) < minimumTransactions {
		t.Fatalf("expected to get transactions number to be greater or equal %d; got: %d", minimumTransactions, len(transactions))
	}
}

func TransferFundsFromActorToPaymail(ctx context.Context, t *testing.T, actorService *ActorService, paymail string) {
	const funds = 2

	transaction, err := actorService.TransferFunds(ctx, paymail, funds)
	if err != nil {
		t.Fatalf("LeaderService - TransferFunds failed: could not transfering funds: %d to: %s because: %v", funds, paymail, err)
	}
	if transaction == nil {
		t.Fatalf("expected to get non nil transaction after transfering funds: %d to: %s paymail", funds, paymail)
	}
	if transaction.OutputValue > -1 {
		t.Fatalf("expected to get output value greater or equal -1, got: %d.", transaction.OutputValue)
	}
}

func TransferFundsFromLeaderToActor(ctx context.Context, t *testing.T, leaderService *LeaderService, actorService *ActorService, funds uint64) {
	actorPaymail := actorService.actor.paymail

	transaction, err := leaderService.TransferFunds(ctx, actorService.actor.paymail, funds)
	if err != nil {
		t.Fatalf("LeaderService - TransferFunds failed: could not transfering funds: %d to: %s because: %v", funds, actorPaymail, err)
	}
	if transaction == nil {
		t.Fatalf("expected to get non nil transaction after transfering funds: %d to: %s paymail", funds, actorPaymail)
	}
	if transaction.OutputValue > -1 {
		t.Fatalf("expected to get output value greater or equal -1, got: %d.", transaction.OutputValue)
	}

	balance, err := leaderService.Balance(ctx)
	if err != nil {
		t.Fatalf("LeaderService - Balance failed: could not check balance status: %v", err)
	}
	if balance < 1 {
		t.Fatalf("expected to get balance value greater or equal 1, got: %d", balance)
	}

	actorTransactions, err := actorService.Transactions(ctx)
	if err != nil {
		t.Fatalf("ActorService - Transactions failed: coult not get transactions: %v", err)
	}
	if actorTransactions == nil {
		t.Fatalf("expected to get non nil transactions slice for %s paymail", actorPaymail)
	}
	if len(actorTransactions) < 1 {
		t.Fatalf("expected to get transactions number to be greater or equal 1; got: %d", len(actorTransactions))
	}
}

func initTestServices(ctx context.Context, url, xPriv, alias string, t *testing.T) (*LeaderService, *ActorService) {
	t.Helper()

	leaderService, err := NewLeaderService(url, xPriv)
	if err != nil {
		t.Fatalf("Setup failed: could not initialize leader service: %v", err)
	}
	actor, err := leaderService.CreateActor(ctx, alias)
	if err != nil {
		t.Fatalf("Setup failed: could not create actor: %v", err)
	}
	actorService, err := NewActorService(leaderService.envURL, actor)
	if err != nil {
		t.Fatalf("Setup failed: count not initialize actor service: %v", err)
	}

	return leaderService, actorService
}
