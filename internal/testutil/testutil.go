package testutil

import (
	"encoding/json"
	"testing"
	"time"

	client "github.com/bitcoin-sv/spv-wallet-go-client"
	"github.com/jarcoal/httpmock"
)

const TestXPriv = "xprv9s21ZrQH143K3fqNnUmXmgfT9ToMtiq5cuKsVBG4E5UqVh4psHDY2XKsEfZKuV4FSZcPS9CYgEQiLUpW2xmHqHFyp23SvTkTCE153cCdwaj"
const TestAPIAddr = "http://localhost:3003"

func NewMockSPVClient(t *testing.T) (*client.Client, *httpmock.MockTransport) {
	t.Helper()
	transport := httpmock.NewMockTransport()
	cfg := client.Config{
		Addr:      TestAPIAddr,
		Timeout:   5 * time.Second,
		Transport: transport,
	}

	spv, err := client.NewWithXPriv(cfg, TestXPriv)
	if err != nil {
		t.Fatalf("test helper - spv wallet client with xpriv: %s", err)
	}

	return spv, transport
}

func MarshalToString(t *testing.T, actual any) string {
	t.Helper()
	bb, err := json.Marshal(actual)
	if err != nil {
		t.Fatalf("test helper - failed to marshal actual response")
	}

	return string(bb)
}
