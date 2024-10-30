package configs

import (
	"context"
	"fmt"

	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/go-resty/resty/v2"
)

const route = "api/v1/configs"

type API struct {
	addr string
	cli  *resty.Client
}

func (a *API) SharedConfig(ctx context.Context) (*response.SharedConfig, error) {
	var spvErr models.SPVError
	var result response.SharedConfig

	resp, err := a.cli.
		R().
		SetContext(ctx).
		SetResult(&result).
		SetError(&spvErr).
		Get(a.addr + "/shared")
	if err != nil {
		return nil, fmt.Errorf("HTTP error: %w", err)
	}
	if resp.IsError() {
		return nil, spvErr
	}
	return &result, nil
}

func NewAPI(addr string, cli *resty.Client) *API {
	return &API{
		addr: addr + "/" + route,
		cli:  cli,
	}
}
