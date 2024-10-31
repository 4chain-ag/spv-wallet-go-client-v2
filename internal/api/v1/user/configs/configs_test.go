package configs_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-go-client/internal/testutil"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestConfigsAPI_SharedConfig_API_ResponseWithStatusCode200(t *testing.T) {
	defer httpmock.DeactivateAndReset()

	// given:
	URL := testutil.TestAPIAddr + "/api/v1/configs/shared"
	spv, transport := testutil.NewMockSPVClient(t)
	file := httpmock.File("configstest/response_200_status_code.json")
	transport.RegisterResponder(http.MethodGet, URL, httpmock.NewJsonResponderOrPanic(http.StatusOK, file))

	// when:
	ctx := context.Background()
	res, err := spv.SharedConfig(ctx)

	// then:
	require.NoError(t, err)
	require.JSONEq(t, file.String(), testutil.MarshalToString(t, res))
}
