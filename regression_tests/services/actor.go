package services

import (
	"fmt"

	"github.com/bitcoin-sv/spv-wallet-go-client/walletkeys"
)

type Actor struct {
	alias   string
	xPriv   string
	xPub    string
	paymail string
}

func (a *Actor) SetPaymail(s string) {
	a.paymail = s
}

func (a *Actor) Paymail() string { return a.paymail }

func NewActor(alias string) (*Actor, error) {
	keys, err := walletkeys.RandomKeys()
	if err != nil {
		return nil, fmt.Errorf("RandomKeys failed: %w", err)
	}

	return &Actor{
		alias: alias,
		xPriv: keys.XPriv(),
		xPub:  keys.XPub(),
	}, nil
}
