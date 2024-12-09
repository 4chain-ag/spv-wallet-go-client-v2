package config_test

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

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

func TestLoadOrDefaultConfig(t *testing.T) {
	// Setup: Create temporary YAML files for testing
	validYAML := `
addr: "http://api.example.com"
timeout: 30s
`
	invalidYAML := `
addr: "http://api.example.com"
timeout: "not-a-number"
`
	missingFile := "nonexistent.yaml"
	validFile := "valid_config.yaml"
	invalidFile := "invalid_config.yaml"

	// Write valid YAML file
	err := os.WriteFile(validFile, []byte(validYAML), 0644)
	require.NoError(t, err)
	defer os.Remove(validFile) // Clean up after test

	// Write invalid YAML file
	err = os.WriteFile(invalidFile, []byte(invalidYAML), 0644)
	require.NoError(t, err)
	defer os.Remove(invalidFile) // Clean up after test

	tests := []struct {
		name        string
		filePath    string
		expectError bool
		expected    config.Config
	}{
		{
			name:        "Valid YAML File",
			filePath:    validFile,
			expectError: false,
			expected: config.Config{
				Addr:    "http://api.example.com",
				Timeout: 30 * time.Second,
			},
		},
		{
			name:        "Invalid YAML File",
			filePath:    invalidFile,
			expectError: true,
			expected:    config.NewConfig(),
		},
		{
			name:        "Missing YAML File",
			filePath:    missingFile,
			expectError: true,
			expected:    config.NewConfig(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg := config.LoadOrDefaultConfig(test.filePath)
			if test.expectError {
				// If an error is expected, ensure default values are set
				require.Equal(t, "http://localhost:3003", cfg.Addr)
				require.Equal(t, 1*time.Minute, cfg.Timeout)
				require.NotNil(t, cfg.Transport)
			} else {
				// If no error is expected, compare expected values
				require.Equal(t, test.expected.Addr, cfg.Addr)
				require.Equal(t, test.expected.Timeout, cfg.Timeout)
				require.NotNil(t, cfg.Transport)
			}
		})
	}
}
