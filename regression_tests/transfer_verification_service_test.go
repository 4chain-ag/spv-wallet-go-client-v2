//go:build regression
// +build regression

package regressiontests

import (
	"context"
	"os"
	"testing"
)

// transferVerifier defines the methods required for verifying a transfer between a sender and a recipient.
// It includes methods for performing the transfer verification (do), retrieving the amount of funds to be transferred (funds),
// and obtaining the paymail addresses of the sender and recipient (senderPaymail and recipientPaymail).
type transferVerifier interface {
	do(ctx context.Context) error
	funds() uint64
	senderPaymail() string
	recipientPaymail() string
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
	leader1, actor1 := initTestServices(ctx, t, alias1, &leaderServiceConfig{
		envURL:     lookupEnvOrDefault(clientOneURL, ""),
		envXPriv:   lookupEnvOrDefault(clientOneLeaderXPriv, ""),
		adminXPriv: admin,
	})
	leader2, actor2 := initTestServices(ctx, t, alias2, &leaderServiceConfig{
		envURL:     lookupEnvOrDefault(clientTwoURL, ""),
		envXPriv:   lookupEnvOrDefault(clientTwoLeaderXPriv, ""),
		adminXPriv: admin,
	})

	t.Cleanup(func() {
		err := leader1.removeActors(ctx)
		if err != nil {
			t.Errorf("Cleanup failed: could not remove actors by leader1 service: %v", err)
		}
	})

	t.Cleanup(func() {
		err := leader2.removeActors(ctx)
		if err != nil {
			t.Errorf("Cleanup failed: could not remove actors by leader2 service: %v", err)
		}
	})

	transfers := []transferVerifier{
		&transferVerficiationService{
			sender:        leader2,
			recipient:     actor1,
			transferFunds: 3,
		},
		&transferVerficiationService{
			sender:        leader1,
			recipient:     actor2,
			transferFunds: 2,
		},
		&transferVerficiationService{
			sender:        actor1,
			recipient:     actor2,
			transferFunds: 2,
		},
	}

	for i, transfer := range transfers {
		err := transfer.do(ctx)
		if err != nil {
			const format = "Transfer no. %d failed: could not transfer: %d funds from Sender [%s] to Recipient [%s]: %v"
			t.Fatalf(format, i+1, transfer.funds(), transfer.senderPaymail(), transfer.recipientPaymail(), err)
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

func initTestServices(ctx context.Context, t *testing.T, alias string, cfg *leaderServiceConfig) (*leaderService, *actorService) {
	t.Helper()

	service1, err := newLeaderService(cfg)
	if err != nil {
		t.Fatalf("Setup failed: could not initialize leader service [url: %s]: %v", cfg.envURL, err)
	}
	actor, err := service1.createActor(ctx, alias)
	if err != nil {
		t.Fatalf("Setup failed: could not create actor [alias: %s]: %v", alias, err)
	}
	service2, err := newActorService(cfg.envURL, actor)
	if err != nil {
		t.Fatalf("Setup failed: initialize actor service [alias: %s]: %v", alias, err)
	}
	return service1, service2
}
