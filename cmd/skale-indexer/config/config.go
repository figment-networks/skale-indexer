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
	DatabaseURL   string `json:"database_url" envconfig:"DATABASE_URL" required:"true"`
	AppEnv        string `json:"app_env" envconfig:"APP_ENV" default:"development"`
	EnableScraper bool   `json:"enable_scraper" envconfig:"ENABLE_SCRAPER" default:"false"`

	Address  string `json:"address" envconfig:"ADDRESS" default:"0.0.0.0"`
	Port     string `json:"port" envconfig:"PORT" default:"3000"`
	HTTPPort string `json:"http_port" envconfig:"HTTP_PORT" default:"8087"`

	EthereumAddress           string `json:"ethereum_address" envconfig:"ETHEREUM_ADDRESS" default:"http://0.0.0.0:8545"`
	SkaleABIDir               string `json:"abi_dir" envconfig:"ABI_DIR" default:"./abi"`
	LowerThresholdForBackward uint64 `json:"lower_threshold_for_backward" envconfig:"LOWER_THRESHOLD_FOR_BACKWARD" default:"0"`

	// Rollbar
	RollbarAccessToken string `json:"rollbar_access_token" envconfig:"ROLLBAR_ACCESS_TOKEN"`
	RollbarServerRoot  string `json:"rollbar_server_root" envconfig:"ROLLBAR_SERVER_ROOT" default:"github.com/figment-networks/skale-indexer"`

	EthereumNodeType string `json:"ethereum_node_type" envconfig:"ETHEREUM_NODE_TYPE" default:"archive"`

	EthereumSmallestBlockNumber uint64 `json:"smallest_block_number" envconfig:"ETHEREUM_SMALLEST_BLOCK_NUMBER"`
	EthereumSmallestTime        uint64 `json:"smallest_block_time" envconfig:"ETHEREUM_SMALLEST_BLOCK_TIME"`

	MaxHeightsPerRequest uint64 `json:"max_heights_per_request" envconfig:"MAX_HEIGHTS_PER_REQUEST" default:"1000"`
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
