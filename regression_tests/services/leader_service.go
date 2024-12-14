package services

import (
	"context"
	"fmt"

	wallet "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/config"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

type LeaderServiceConfig struct {
	EnvURL     string
	EnvXPriv   string
	AdminXPriv string
	Domain     string
	Name       string
}

type LeaderService struct {
	cfg      *LeaderServiceConfig
	adminAPI *wallet.AdminAPI
	userAPI  *wallet.UserAPI
	actors   []*Actor
}

func (l *LeaderService) Name() string { return l.cfg.Name }

func (l *LeaderService) CreatePaymail(alias string) string { return alias + "@" + l.cfg.Domain }

func (l *LeaderService) FirstActor() *Actor {
	if len(l.actors) > 0 {
		return l.actors[0]
	}
	return nil
}

func (l *LeaderService) Balance(ctx context.Context) (uint64, error) {
	xPub, err := l.userAPI.XPub(ctx)
	if err != nil {
		return 0, fmt.Errorf("XPub failed: %w", err)
	}
	return xPub.CurrentBalance, nil
}

func (l *LeaderService) CreateActor(ctx context.Context, alias string) (*Actor, error) {
	actor, err := NewActor(alias)
	if err != nil {
		return nil, fmt.Errorf("NewActor failed: %w", err)
	}

	_, err = l.adminAPI.CreateXPub(ctx, &commands.CreateUserXpub{XPub: actor.xPub})
	if err != nil {
		return nil, fmt.Errorf("CreateXPub failed: %w", err)
	}

	actor.SetPaymail(l.CreatePaymail(alias))
	_, err = l.adminAPI.CreatePaymail(ctx, &commands.CreatePaymail{
		Key:        actor.xPub,
		Address:    actor.paymail,
		PublicName: "Regression tests",
	})
	if err != nil {
		return nil, fmt.Errorf("CreatePaymail failed: %w", err)
	}

	l.actors = append(l.actors, actor)
	return actor, nil
}

func (l *LeaderService) RemoveActors(ctx context.Context) error {
	for _, a := range l.actors {
		err := l.adminAPI.DeletePaymail(ctx, a.paymail)
		if err != nil {
			return fmt.Errorf("DeletePaymail failed: %w", err)
		}
	}
	return nil
}

func (l *LeaderService) TransferFunds(ctx context.Context, paymail string, satoshis uint64) (*response.Transaction, error) {
	balance, err := l.Balance(ctx)
	if err != nil {
		return nil, fmt.Errorf("Balance check failed: %w", err)
	}
	if balance < satoshis {
		const format = "Failed to transfer %d satoshis to paymail: %s. Current leader balance: %d"
		return nil, fmt.Errorf(format, satoshis, paymail, balance)
	}

	transaction, err := l.userAPI.SendToRecipients(ctx, &commands.SendToRecipients{
		Recipients: []*commands.Recipients{
			{
				To: paymail, Satoshis: satoshis,
			},
		},
		Metadata: map[string]any{"description": "regression-test"},
	})
	if err != nil {
		return nil, fmt.Errorf("SendToRecipients failed: %w", err)
	}
	return transaction, nil

}

func NewLeaderService(cfg *LeaderServiceConfig) (*LeaderService, error) {
	walletCfg := config.New(config.WithAddr(cfg.EnvURL))
	adminAPI, err := wallet.NewAdminAPIWithXPriv(walletCfg, cfg.AdminXPriv)
	if err != nil {
		return nil, fmt.Errorf("NewAdminAPIWithXPriv failed: %w", err)
	}
	userAPI, err := wallet.NewUserAPIWithXPriv(walletCfg, cfg.EnvXPriv)
	if err != nil {
		return nil, fmt.Errorf("NewUserAPIWithXPriv failed: %w", err)
	}
	sharedCfg, err := userAPI.SharedConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("SharedConfig failed: %w", err)
	}
	if len(sharedCfg.PaymailDomains) != 1 {
		return nil, fmt.Errorf("Expected to have single paymail domain, got: %d", len(sharedCfg.PaymailDomains))
	}

	cfg.Domain = sharedCfg.PaymailDomains[0]
	cfg.Name = "Leader@" + cfg.Domain

	return &LeaderService{
		cfg:      cfg,
		adminAPI: adminAPI,
		userAPI:  userAPI,
	}, nil
}
