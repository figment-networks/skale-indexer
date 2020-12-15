package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/kelseyhightower/envconfig"
)

var (
	Name      = "skale-indexer"
	Version   string
	GitSHA    string
	Timestamp string
)

const (
	modeDevelopment = "development"
	modeProduction  = "production"
)

type Config struct {
	DatabaseURL string `json:"database_url" envconfig:"DATABASE_URL" required:"true"`
	AppEnv      string `json:"app_env" envconfig:"APP_ENV" default:"development"`

	Address  string `json:"address" envconfig:"ADDRESS" default:"0.0.0.0"`
	Port     string `json:"port" envconfig:"PORT" default:"3000"`
	HTTPPort string `json:"http_port" envconfig:"HTTP_PORT" default:"8087"`

	EthereumAddress string `json:"ethereum_address" envconfig:"ETHEREUM_ADDRESS" default:"http://0.0.0.0:8545"`
	SkaleABIDir     string `json:"abi_dir" envconfig:"ABI_DIR" default:"./abi"`

	// Rollbar
	RollbarAccessToken string `json:"rollbar_access_token" envconfig:"ROLLBAR_ACCESS_TOKEN"`
	RollbarServerRoot  string `json:"rollbar_server_root" envconfig:"ROLLBAR_SERVER_ROOT" default:"github.com/figment-networks/skale-indexer"`
}

// IdentityString returns the full app version string
func IdentityString() string {

	t, err := strconv.Atoi(Timestamp)
	timestamp := Timestamp
	if err == nil {
		timestamp = time.Unix(int64(t), 0).String()
	}
	return fmt.Sprintf(
		"%s %s (git: %s) - built at %s",
		Name,
		Version,
		GitSHA,
		timestamp,
	)
}

// FromFile reads the config from a file
func FromFile(path string, config *Config) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, config)
}

// FromEnv reads the config from environment variables
func FromEnv(config *Config) error {
	return envconfig.Process("", config)
}
