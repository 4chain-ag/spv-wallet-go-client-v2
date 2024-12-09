package config

import (
	"log"
	"net/http"
	"net/url"
	"time"

	goclienterr "github.com/bitcoin-sv/spv-wallet-go-client/errors"
)

// Config holds configuration settings for establishing a connection and handling
// request details in the application.
type Config struct {
	Addr      string            // The base address of the SPV Wallet API.
	Timeout   time.Duration     // The HTTP requests timeout duration.
	Transport http.RoundTripper // Custom HTTP transport, allowing optional customization of the HTTP client behavior.
}

// NewConfig creates a new Config instance with optional customizations.
func NewConfig(options ...Option) Config {
	cfg := Config{}
	for _, opt := range options {
		opt(&cfg)
	}
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Error creating configuration: %v", err)
	}

	cfg.setDefaultValues()
	return cfg
}

// Validate checks the configuration for invalid or missing values.
func (cfg *Config) Validate() error {
	if cfg.Addr == "" {
		return goclienterr.ErrConfigValidationMissingAddress
	}

	if _, err := url.ParseRequestURI(cfg.Addr); err != nil {
		return goclienterr.ErrConfigValidationInvalidAddress
	}

	if cfg.Timeout <= 0 {
		return goclienterr.ErrConfigValidationInvalidTimeout
	}

	if cfg.Transport == nil {
		return goclienterr.ErrConfigValidationInvalidTransport
	}

	return nil
}
