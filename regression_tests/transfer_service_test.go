//go:build regression
// +build regression

package regressiontests

import (
	"context"
	"os"
	"testing"
)

func TestTransferService_Do(t *testing.T) {
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
	leader1, actor1 := initTestServices(ctx, t, alias1, &LeaderServiceConfig{
		EnvURL:     lookupEnvOrDefault(clientOneURL, ""),
		EnvXPriv:   lookupEnvOrDefault(clientOneLeaderXPriv, ""),
		AdminXPriv: admin,
	})
	leader2, actor2 := initTestServices(ctx, t, alias2, &LeaderServiceConfig{
		EnvURL:     lookupEnvOrDefault(clientTwoURL, ""),
		EnvXPriv:   lookupEnvOrDefault(clientTwoLeaderXPriv, ""),
		AdminXPriv: admin,
	})

	t.Cleanup(func() {
		err := leader1.RemoveActors(ctx)
		if err != nil {
			t.Errorf("Cleanup failed: could not remove actors by leader1 service: %v", err)
		}
	})

	t.Cleanup(func() {
		err := leader2.RemoveActors(ctx)
		if err != nil {
			t.Errorf("Cleanup failed: could not remove actors by leader2 service: %v", err)
		}
	})

	firstTransfer := &TransferService{
		Sender:        leader2,
		Recipient:     actor1,
		TransferFunds: 3,
	}
	err := firstTransfer.Do(ctx)
	if err != nil {
		t.Errorf("Transfer no. 1 failed: could not transfer: %d funds from leader2 [url: %s] to actor1 [alias: %s]: %v", firstTransfer.TransferFunds, leader2.cfg.EnvURL, actor1.actor.alias, err)
	}

	secondTransfer := &TransferService{
		Sender:        leader1,
		Recipient:     actor2,
		TransferFunds: 2,
	}
	err = secondTransfer.Do(ctx)
	if err != nil {
		t.Errorf("Transfer no. 2 failed: could not transfer: %d funds from leader1 [url: %s] to actor2 [alias: %s]: %v", secondTransfer.TransferFunds, leader2.cfg.EnvURL, actor2.actor.alias, err)
	}

	thirdTransfer := &TransferService{
		Sender:        actor1,
		Recipient:     actor2,
		TransferFunds: 2,
	}
	err = thirdTransfer.Do(ctx)
	if err != nil {
		t.Errorf("Transfer no. 3 failed: could not transfer: %d funds from actor1 [alias: %s] to actor2 [alias: %s]: %v", thirdTransfer.TransferFunds, actor1.actor.alias, actor2.actor.alias, err)
	}
}

func lookupEnvOrDefault(env string, s string) string {
	v, ok := os.LookupEnv(env)
	if ok {
		return v
	}
	return s
}

func initTestServices(ctx context.Context, t *testing.T, alias string, cfg *LeaderServiceConfig) (*LeaderService, *ActorService) {
	t.Helper()

	service1, err := NewLeaderService(cfg)
	if err != nil {
		t.Fatalf("Setup failed: could not initialize leader service [url: %s]: %v", cfg.EnvURL, err)
	}
	actor, err := service1.CreateActor(ctx, alias)
	if err != nil {
		t.Fatalf("Setup failed: could not create actor [alias: %s]: %v", alias, err)
	}
	service2, err := NewActorService(cfg.EnvURL, actor)
	if err != nil {
		t.Fatalf("Setup failed: initialize actor service [alias: %s]: %v", alias, err)
	}
	return service1, service2
}
