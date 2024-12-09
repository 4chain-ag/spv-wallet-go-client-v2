package exampleutil

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/bitcoin-sv/spv-wallet-go-client/config"
)

// NewDefaultConfig returns a new instance of the default example configuration.
func NewDefaultConfig() config.Config {
	return config.NewConfig()
}

// LoadConfigFromFile loads a configuration from a file.
func LoadConfigFromFile(filePath string) config.Config {
	return config.LoadOrDefaultConfig(filePath)
}

// NewCustomConfig allows creating a custom configuration with optional parameters.
func NewCustomConfig(addr string, timeout time.Duration, transport *http.Transport) config.Config {
	options := []config.Option{}

	if addr != "" {
		options = append(options, config.WithAddr(addr))
	}

	if timeout > 0 {
		options = append(options, config.WithTimeout(timeout))
	}

	if transport != nil {
		options = append(options, config.WithTransport(transport))
	}

	return config.NewConfig(options...)
}

func Print(s string, a any) {
	fmt.Println(strings.Repeat("~", 100))
	fmt.Println(s)
	fmt.Println(strings.Repeat("~", 100))
	res, err := json.MarshalIndent(a, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(res))
}

func RandomPaymail() string {
	seed := time.Now().UnixNano()
	n := rand.New(rand.NewSource(seed)).Intn(500)
	addr := fmt.Sprintf("john.doe.%dtest@4chain.test.com", n)
	return addr
}
