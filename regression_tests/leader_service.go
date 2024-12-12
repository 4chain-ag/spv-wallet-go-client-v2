package regressiontests

import (
	"context"
	"fmt"

	wallet "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/config"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/querybuilders"
	"github.com/bitcoin-sv/spv-wallet-go-client/walletkeys"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

type Actor struct {
	xPriv   string
	xPub    string
	paymail string
}

type LeaderService struct {
	adminAPI   *wallet.AdminAPI
	userAPI    *wallet.UserAPI
	actors     []*Actor
	name       string
	domain     string
	adminXPriv string
	envURL     string
	envXPriv   string
}

func (l *LeaderService) Balance(ctx context.Context) (uint64, error) {
	xPub, err := l.userAPI.XPub(ctx)
	if err != nil {
		return 0, fmt.Errorf("UserAPI - XPub failed: %v", err)
	}
	return xPub.CurrentBalance, nil
}

func (l *LeaderService) TransferFunds(ctx context.Context, paymail string, ammount uint64) (*response.Transaction, error) {
	balance, err := l.Balance(ctx)
	if err != nil {
		return nil, fmt.Errorf("Balance check failed: %v", err)
	}
	if balance < ammount {
		const format = "Failed to transfer %d satoshis to paymail: %s. Current leader balance: %d"
		return nil, fmt.Errorf(format, ammount, paymail, balance)
	}

	transaction, err := l.userAPI.SendToRecipients(ctx, &commands.SendToRecipients{
		Recipients: []*commands.Recipients{
			{
				To: paymail, Satoshis: ammount,
			},
		},
		Metadata: map[string]any{"description": "regression-test"},
	})
	if err != nil {
		return nil, fmt.Errorf("UserAPI - SendToRecipients failed: %v", err)
	}
	return transaction, nil
}

func (l *LeaderService) RemoveActors(ctx context.Context) error {
	for _, a := range l.actors {
		err := l.adminAPI.DeletePaymail(ctx, a.paymail)
		if err != nil {
			return fmt.Errorf("AdminAPI - DeletePaymail failed: %v", err)
		}
	}
	return nil
}

func (l *LeaderService) CreateActor(ctx context.Context, alias string) (*Actor, error) {
	keys, err := walletkeys.RandomKeys()
	if err != nil {
		return nil, fmt.Errorf("RandomKeys failed: %v", err)
	}

	xPub := keys.XPub()
	_, err = l.adminAPI.CreateXPub(ctx, &commands.CreateUserXpub{
		Metadata: make(querybuilders.Metadata),
		XPub:     xPub,
	})
	if err != nil {
		return nil, fmt.Errorf("AdminAPI - CreateXPub failed: %v", err)
	}

	paymail := alias + "@" + l.domain
	_, err = l.adminAPI.CreatePaymail(ctx, &commands.CreatePaymail{
		Metadata:   make(querybuilders.Metadata),
		Key:        xPub,
		Address:    paymail,
		PublicName: "Regression tests",
	})
	if err != nil {
		return nil, fmt.Errorf("AdminAPI - CreatePaymail failed: %v", err)
	}

	actor := Actor{
		xPriv:   keys.XPriv(),
		xPub:    xPub,
		paymail: paymail,
	}
	l.actors = append(l.actors, &actor)
	return &actor, nil
}

func NewLeaderService(urlEnv, xPrivEnv string) (*LeaderService, error) {
	const adminXPriv = "xprv9s21ZrQH143K3CbJXirfrtpLvhT3Vgusdo8coBritQ3rcS7Jy7sxWhatuxG5h2y1Cqj8FKmPp69536gmjYRpfga2MJdsGyBsnB12E19CESK"
	envXPriv := setEnvXPrivOrDefault(xPrivEnv)
	envURL := setEnvURLOrDefault(urlEnv)
	cfg := config.New(config.WithAddr(envURL))

	adminAPI, err := wallet.NewAdminAPIWithXPriv(cfg, adminXPriv)
	if err != nil {
		return nil, err
	}
	userAPI, err := wallet.NewUserAPIWithXPriv(cfg, envXPriv)
	if err != nil {
		return nil, err
	}
	sharedCfg, err := userAPI.SharedConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("UserAPI - SharedConfig failed: %v", err)
	}
	if len(sharedCfg.PaymailDomains) != 1 {
		return nil, fmt.Errorf("expected to have single paymail domain, got: %d", len(sharedCfg.PaymailDomains))
	}

	domain := sharedCfg.PaymailDomains[0]
	return &LeaderService{
		adminAPI:   adminAPI,
		userAPI:    userAPI,
		adminXPriv: adminXPriv,
		name:       "Leader@" + domain,
		domain:     domain,
		envURL:     envURL,
		envXPriv:   envXPriv,
	}, nil
}
