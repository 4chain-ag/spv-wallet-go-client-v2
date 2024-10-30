package configs

import (
	"context"
	"fmt"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/restyutil"
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
	var result response.SharedConfig
	response := restyutil.ResponseAdapter(func() (*resty.Response, error) {
		return a.cli.
			R().
			SetContext(ctx).
			SetResult(&result).
			SetError(&models.SPVError{}).
			Get(a.addr + "/shared")
	})
	err := response.HandleErr()
	if err != nil {
		return nil, fmt.Errorf("HTTP response failure: %w", err)
	}
	return &result, nil
}

func NewAPI(addr string, cli *resty.Client) *API {
	return &API{
		addr: addr + "/" + route,
		cli:  cli,
	}
}
