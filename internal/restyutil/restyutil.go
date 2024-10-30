package restyutil

import (
	"errors"
	"fmt"

	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/go-resty/resty/v2"
)

type ResponseAdapter func() (*resty.Response, error)

func (r ResponseAdapter) HandleErr() error {
	res, err := r()
	if err != nil {
		return fmt.Errorf("HTTP error: %w", err)
	}
	if res.IsSuccess() {
		return nil
	}
	if v, ok := res.Error().(*models.SPVError); ok && len(v.Code) > 0 {
		return v
	}
	return ErrUnrecognizedAPIResponse
}

var ErrUnrecognizedAPIResponse = errors.New("unrecognized response from API")
