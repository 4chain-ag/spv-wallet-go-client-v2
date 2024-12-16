package regressiontests

import (
	"context"
	"fmt"

	wallet "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/bitcoin-sv/spv-wallet-go-client/commands"
	"github.com/bitcoin-sv/spv-wallet-go-client/config"
	"github.com/bitcoin-sv/spv-wallet-go-client/walletkeys"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

type Actor struct {
	alias   string
	xPriv   string
	xPub    string
	paymail string
}

func NewActor(alias, domain string) (*Actor, error) {
	keys, err := walletkeys.RandomKeys()
	if err != nil {
		return nil, fmt.Errorf("")
	}
	return &Actor{
		alias:   alias,
		xPriv:   keys.XPriv(),
		xPub:    keys.XPub(),
		paymail: alias + "@" + domain,
	}, nil
}

type ActorService struct {
	userAPI *wallet.UserAPI
	actor   *Actor
}

func (a *ActorService) Paymail() string { return a.actor.paymail }

func (a *ActorService) Balance(ctx context.Context) (uint64, error) {
	xpub, err := a.userAPI.XPub(ctx)
	if err != nil {
		return 0, fmt.Errorf("could not fetch xPub to retrieve current balance: %w", err)
	}
	return xpub.CurrentBalance, nil
}

func (a *ActorService) Transactions(ctx context.Context) ([]*response.Transaction, error) {
	page, err := a.userAPI.Transactions(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not fetch transactions: %w", err)
	}
	return page.Content, nil
}

func (a *ActorService) TransferFunds(ctx context.Context, paymail string, funds uint64) (*response.Transaction, error) {
	balance, err := a.Balance(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not fetch balance: %w", err)
	}
	if balance < funds {
		return nil, fmt.Errorf("insufficient balance: %d to spend: %d in transaction", balance, funds)
	}

	recipient := commands.Recipients{To: paymail, Satoshis: funds}
	transaction, err := a.userAPI.SendToRecipients(ctx, &commands.SendToRecipients{
		Recipients: []*commands.Recipients{&recipient},
		Metadata:   map[string]any{"description": "regression-test"},
	})
	if err != nil {
		return nil, fmt.Errorf("could not transfer funds: %w", err)
	}
	return transaction, nil
}

func NewActorService(url string, actor *Actor) (*ActorService, error) {
	userAPI, err := wallet.NewUserAPIWithXPriv(config.New(config.WithAddr(url)), actor.xPriv)
	if err != nil {
		return nil, fmt.Errorf("could not initialize user API: %w", err)
	}
	return &ActorService{userAPI: userAPI, actor: actor}, nil
}
