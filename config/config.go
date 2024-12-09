package config

import (
	"net/http"
	"time"
)

// Config holds configuration settings for establishing a connection and handling
// request details in the application.
type Config struct {
	Addr      string            // The base address of the SPV Wallet API.
	Timeout   time.Duration     // The HTTP requests timeout duration.
	Transport http.RoundTripper // Custom HTTP transport, allowing optional customization of the HTTP client behavior.
}

// setDefaultValues assigns default values to fields that are not explicitly set.
func (cfg *Config) setDefaultValues() {
	if cfg.Addr == "" {
		cfg.Addr = "http://localhost:3003"
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 1 * time.Minute
	}
	if cfg.Transport == nil {
		cfg.Transport = http.DefaultTransport
	}
}

// Option defines a function signature for modifying a Config.
type Option func(*Config)

// WithAddr sets the address in the configuration.
func WithAddr(addr string) Option {
	return func(cfg *Config) {
		cfg.Addr = addr
	}
}

// WithTimeout sets the timeout duration in the configuration.
func WithTimeout(timeout time.Duration) Option {
	return func(cfg *Config) {
		cfg.Timeout = timeout
	}
}

// WithTransport sets the HTTP transport in the configuration.
func WithTransport(transport http.RoundTripper) Option {
	return func(cfg *Config) {
		cfg.Transport = transport
	}
}

// NewConfig creates a new Config instance with optional customizations.
func NewConfig(options ...Option) Config {
	cfg := Config{}
	for _, opt := range options {
		opt(&cfg)
	}
	cfg.setDefaultValues()
	return cfg
}
