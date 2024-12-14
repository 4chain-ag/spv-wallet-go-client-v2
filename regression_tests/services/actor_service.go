package services

import (
	"context"
	"fmt"

	wallet "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/config"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

type ActorService struct {
	url     string
	userAPI *wallet.UserAPI
	actor   *Actor
}

func (a *ActorService) Actor() *Actor { return a.actor }

func (a *ActorService) Balance(ctx context.Context) (uint64, error) {
	xPub, err := a.userAPI.XPub(ctx)
	if err != nil {
		return 0, fmt.Errorf("XPub failed: %w", err)
	}
	return xPub.CurrentBalance, nil
}

func (a *ActorService) TransferFunds(ctx context.Context, paymail string, funds uint64) (*response.Transaction, error) {
	balance, err := a.Balance(ctx)
	if err != nil {
		return nil, fmt.Errorf("Balance check failed: %w", err)
	}
	if balance < funds {
		return nil, fmt.Errorf("Failed to transfer %d satoshis to paymail: %s. Current balance: %d.", funds, paymail, balance)
	}

	recipient := commands.Recipients{To: paymail, Satoshis: funds}
	transaction, err := a.userAPI.SendToRecipients(ctx, &commands.SendToRecipients{
		Recipients: []*commands.Recipients{&recipient},
		Metadata:   map[string]any{"description": "regression-test"},
	})
	if err != nil {
		return nil, fmt.Errorf("SendToRecipients failed: %w", err)
	}
	return transaction, nil
}

func (a *ActorService) Transactions(ctx context.Context) ([]*response.Transaction, error) {
	page, err := a.userAPI.Transactions(ctx)
	if err != nil {
		return nil, fmt.Errorf("Transactions failed: %w", err)
	}
	return page.Content, nil
}

func NewActorService(url string, actor *Actor) (*ActorService, error) {
	userAPI, err := wallet.NewUserAPIWithXPriv(config.New(config.WithAddr(url)), actor.xPriv)
	if err != nil {
		return nil, fmt.Errorf("NewUserAPIWithXPriv failed: %w", err)
	}

	return &ActorService{
		url:     url,
		userAPI: userAPI,
		actor:   actor,
	}, nil
}
