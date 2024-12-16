package regressiontests

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
}

type LeaderService struct {
	cfg      *LeaderServiceConfig
	adminAPI *wallet.AdminAPI
	userAPI  *wallet.UserAPI
	actors   []*Actor
	domain   string
}

func (l *LeaderService) Paymail() string { return "Leader@" + l.domain }

func (l *LeaderService) Balance(ctx context.Context) (uint64, error) {
	xpub, err := l.userAPI.XPub(ctx)
	if err != nil {
		return 0, fmt.Errorf("could not fetch xPub to retrieve current balance: %w", err)
	}
	return xpub.CurrentBalance, nil
}

func (l *LeaderService) Transactions(ctx context.Context) ([]*response.Transaction, error) {
	page, err := l.userAPI.Transactions(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not fetch transactions: %w", err)
	}
	return page.Content, nil
}

func (l *LeaderService) TransferFunds(ctx context.Context, paymail string, funds uint64) (*response.Transaction, error) {
	balance, err := l.Balance(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not fetch balance: %w", err)
	}
	if balance < funds {
		return nil, fmt.Errorf("insufficient balance: %d to spend: %d in transaction", balance, funds)
	}

	recipient := commands.Recipients{To: paymail, Satoshis: funds}
	transaction, err := l.userAPI.SendToRecipients(ctx, &commands.SendToRecipients{
		Recipients: []*commands.Recipients{&recipient},
		Metadata:   map[string]any{"description": "regression-test"},
	})
	if err != nil {
		return nil, fmt.Errorf("could not transfer funds: %w", err)
	}
	return transaction, nil
}

func (l *LeaderService) CreateActor(ctx context.Context, alias string) (*Actor, error) {
	actor, err := NewActor(alias, l.domain)
	if err != nil {
		return nil, fmt.Errorf("could not create the actor: %w", err)
	}

	_, err = l.adminAPI.CreateXPub(ctx, &commands.CreateUserXpub{XPub: actor.xPub})
	if err != nil {
		return nil, fmt.Errorf("could not create the xPub: %w", err)
	}
	_, err = l.adminAPI.CreatePaymail(ctx, &commands.CreatePaymail{
		Key:        actor.xPub,
		Address:    actor.paymail,
		PublicName: "Regression tests",
	})
	if err != nil {
		return nil, fmt.Errorf("could not create the paymail: %w", err)
	}

	l.actors = append(l.actors, actor)
	return actor, nil
}

func (l *LeaderService) RemoveActors(ctx context.Context) error {
	for _, a := range l.actors {
		err := l.adminAPI.DeletePaymail(ctx, a.paymail)
		if err != nil {
			return fmt.Errorf("could not delete the paymail: %w", err)
		}
	}
	return nil
}

func NewLeaderService(cfg *LeaderServiceConfig) (*LeaderService, error) {
	walletCfg := config.New(config.WithAddr(cfg.EnvURL))
	adminAPI, err := wallet.NewAdminAPIWithXPriv(walletCfg, cfg.AdminXPriv)
	if err != nil {
		return nil, fmt.Errorf("could not initialize admin API: %w", err)
	}
	userAPI, err := wallet.NewUserAPIWithXPriv(walletCfg, cfg.EnvXPriv)
	if err != nil {
		return nil, fmt.Errorf("could not initialize user API: %w", err)
	}

	sharedCfg, err := userAPI.SharedConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("could not fetch shared config to retrieve paymail domains: %w", err)
	}
	if len(sharedCfg.PaymailDomains) > 1 {
		return nil, fmt.Errorf("expect to have single paymail domain. Got: %d paymail domains", len(sharedCfg.PaymailDomains))
	}
	return &LeaderService{
		cfg:      cfg,
		adminAPI: adminAPI,
		userAPI:  userAPI,
		domain:   sharedCfg.PaymailDomains[0],
	}, nil
}
