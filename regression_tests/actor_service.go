package regressiontests

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
	actor   *Actor
	userAPI *wallet.UserAPI
}

func (a *ActorService) Balance(ctx context.Context) (uint64, error) {
	xPub, err := a.userAPI.XPub(ctx)
	if err != nil {
		return 0, err
	}

	return xPub.CurrentBalance, nil
}

func (a *ActorService) TransferFunds(ctx context.Context, paymail string, ammount uint64) (*response.Transaction, error) {
	balance, err := a.Balance(ctx)
	if err != nil {
		return nil, fmt.Errorf("balance failure: %v", err)
	}
	if balance < ammount {
		const format = "failed to transfer %d satoshis to paymail: %s. Current balance: %d."
		return nil, fmt.Errorf(format, ammount, paymail, balance)
	}

	transaction, err := a.userAPI.SendToRecipients(ctx, &commands.SendToRecipients{
		Recipients: []*commands.Recipients{
			{
				To: paymail, Satoshis: ammount,
			},
		},
		Metadata: map[string]any{"description": "regression-test"},
	})
	if err != nil {
		return nil, fmt.Errorf("send to recipients failure: %v", err)
	}
	return transaction, nil
}

func (a *ActorService) Transactions(ctx context.Context) ([]*response.Transaction, error) {
	transactions, err := a.userAPI.Transactions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve transactions: %v.", err)
	}
	return transactions.Content, nil
}

func NewActorService(url string, actor *Actor) (*ActorService, error) {
	userAPI, err := wallet.NewUserAPIWithXPriv(config.New(config.WithAddr(url)), actor.xPriv)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize user API: %v", err)
	}

	return &ActorService{
		url:     url,
		actor:   actor,
		userAPI: userAPI,
	}, nil
}
