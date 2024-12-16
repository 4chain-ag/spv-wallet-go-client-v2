//go:build regression
// +build regression

package regressiontests

import (
	"context"
	"os"
	"testing"
)

// TransferVerifier defines the methods required for verifying a transfer between a sender and a recipient.
// It includes methods for performing the transfer verification (Do), retrieving the amount of funds to be transferred (Funds),
// and obtaining the paymail addresses of the sender and recipient (SenderPaymail and RecipientPaymail).
type TransferVerifier interface {
	Do(ctx context.Context) error
	Funds() uint64
	SenderPaymail() string
	RecipientPaymail() string
}

func TestTransferVerficiationService_Do(t *testing.T) {
	const (
		clientOneURL         = "CLIENT_ONE_URL"
		clientOneLeaderXPriv = "CLIENT_ONE_LEADER_XPRIV"
		clientTwoURL         = "CLIENT_TWO_URL"
		clientTwoLeaderXPriv = "CLIENT_TWO_LEADER_XPRIV"
	)

	const (
		alias1 = "Actor11RegressionTest"
		alias2 = "Actor22RegressionTest"
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

	transfers := []TransferVerifier{
		&TransferVerficiationService{
			Sender:        leader2,
			Recipient:     actor1,
			TransferFunds: 3,
		},
		&TransferVerficiationService{
			Sender:        leader1,
			Recipient:     actor2,
			TransferFunds: 2,
		},
		&TransferVerficiationService{
			Sender:        actor1,
			Recipient:     actor2,
			TransferFunds: 2,
		},
	}

	for i, transfer := range transfers {
		err := transfer.Do(ctx)
		if err != nil {
			const format = "Transfer no. %d failed: could not transfer: %d funds from Sender [%s] to Recipient [%s]: %v"
			t.Errorf(format, i+1, transfer.Funds(), transfer.SenderPaymail(), transfer.RecipientPaymail(), err)
		}
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
