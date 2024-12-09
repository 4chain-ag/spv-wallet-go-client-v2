package config

import (
	"log"
	"net/http"
	"time"

	"github.com/spf13/viper"
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
		cfg.Addr = "http://localhost:3003" // Default address
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 1 * time.Minute // Default timeout
	}
	if cfg.Transport == nil {
		cfg.Transport = http.DefaultTransport // Default transport
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

// loadConfigFromYAML loads configuration values from a YAML file using Viper.
func loadConfigFromYAML(filePath string) (Config, error) {
	viper.SetConfigFile(filePath)

	// Read the configuration file
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file: %v", err)
		return NewConfig(), err
	}

	// Set default values
	viper.SetDefault("addr", "http://localhost:3003")
	viper.SetDefault("timeout", 60) // Timeout in seconds

	// Unmarshal into Config struct
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Printf("Error unmarshaling config: %v", err)
		return NewConfig(), err
	}

	// Convert timeout from seconds to time.Duration
	if cfg.Timeout == 0 {
		cfg.Timeout = time.Duration(viper.GetInt("timeout")) * time.Second
	}

	// Set default values for any missing fields
	cfg.setDefaultValues()

	return cfg, nil
}

/*
example yaml file:
---
addr: "http://api.example.com"
timeout: 30
*/

// LoadOrDefaultConfig attempts to load configuration from a YAML file.
// If the file does not exist or an error occurs, it falls back to NewConfig.
func LoadOrDefaultConfig(filePath string) Config {
	cfg, err := loadConfigFromYAML(filePath)
	if err != nil {
		log.Printf("loading default config: %v", err)
		// Fall back to default configuration
		return NewConfig()
	}
	return cfg
}
