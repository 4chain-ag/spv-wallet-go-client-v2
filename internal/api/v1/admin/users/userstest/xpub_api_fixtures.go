package userstest

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/models"
)

func Ptr[T any](value T) *T {
	return &value
}

func NewBadRequestSPVError() *models.SPVError {
	return &models.SPVError{
		Message:    http.StatusText(http.StatusBadRequest),
		StatusCode: http.StatusBadRequest,
		Code:       "invalid-data-format",
	}
}
