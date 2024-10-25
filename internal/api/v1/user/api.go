package user

import (
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/api/v1/user/configurations"
	"github.com/bitcoin-sv/spv-wallet-go-client/internal/httpx"
)

const (
	v1 = "api/v1/"
)

// API aggregates the user-related API endpoint groups exposed by the SPV Wallet service.
//
// This struct serves as an abstraction layer that simplifies interactions with the user APIs.
// It enables developers to easily perform operations related to specific API endpoint groups,
// such as managing user configurations and accessing related functionalities. By utilizing this
// API struct, developers can streamline their code and focus on high-level operations without
// needing to manage the underlying API details directly.
type API struct {
	*configurations.API
}

// NewAPI initializes a new API instance for user-related operations.
//
// This factory function sets up the user API client with the specified base address
// and HTTP client. It constructs the API endpoint path by appending "/api/v1/"
// followed by the specified domain (e.g., "configs") to the provided address.
// This allows for easy access to user-related endpoints, enabling operations
// such as retrieving or updating user configurations. The resulting user API
// instance can be used to make requests to various user-related endpoints.
func NewAPI(addr string, h *httpx.Resty) *API {
	api := API{
		API: &configurations.API{
			Path: addr + v1 + "/configs",
			HTTP: h,
		},
	}
	return &api
}
