package config_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/bitcoin-sv/spv-wallet-go-client/config"
)

func TestNewConfig(t *testing.T) {
	transport := &http.Transport{
		MaxIdleConns:        10,
		MaxIdleConnsPerHost: 10,
		Proxy:               http.ProxyFromEnvironment,
	}

	tests := []struct {
		name     string
		options  []config.Option
		expected config.Config
	}{
		{
			name:    "All defaults",
			options: nil,
			expected: config.Config{
				Addr:      "http://localhost:3003",
				Timeout:   1 * time.Minute,
				Transport: http.DefaultTransport,
			},
		},
		{
			name: "Partial customization",
			options: []config.Option{
				config.WithAddr("http://api.example.com"),
			},
			expected: config.Config{
				Addr:      "http://api.example.com",
				Timeout:   1 * time.Minute,
				Transport: http.DefaultTransport,
			},
		},
		{
			name: "Full customization",
			options: []config.Option{
				config.WithAddr("http://custom.example.com"),
				config.WithTimeout(2 * time.Minute),
				config.WithTransport(transport),
			},
			expected: config.Config{
				Addr:      "http://custom.example.com",
				Timeout:   2 * time.Minute,
				Transport: transport,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg := config.NewConfig(test.options...)
			if cfg != test.expected {
				t.Errorf("Expected %+v, got %+v", test.expected, cfg)
			}
		})
	}
}
