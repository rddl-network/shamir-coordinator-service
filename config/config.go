package config

import "sync"

const DefaultConfigTemplate = `
SERVICE_HOST="{{ .ServiceHost }}"
SERVICE_PORT={{ .ServicePort }}
SHARE_HOLDER_LIST="{{ .ShareHolderList }}"
CERTS_PATH="{{ .CertsPath }}"
WALLET_NAME="{{ .WalletName }}"
ELEMENTS_RPC_URL="{{ .ElementsRpcUrl }}"
ASSET_ID="{{ .AssetID }}"
SHAMIR_THRESHOLD={{ .ShamirThreshold }}
SHAMIR_SHARES={{ .ShamirShares }}
TEST_MODE={{ .TestMode }}
`

// Config defines TA's top level configuration
type Config struct {
	ServiceBind     string `json:"service-bind"        mapstructure:"service-bind"`
	ServicePort     int    `json:"service-port"        mapstructure:"service-port"`
	ShareHolderList string `json:"share-holder-list"   mapstructure:"share-holder-list"`
	CertsPath       string `json:"certs-path"          mapstructure:"certs-path"`
	WalletName      string `json:"wallet-name"         mapstructure:"wallet-name"`
	ElementsRpcUrl  string `json:"elements-rpc-url"    mapstructure:"elements-rpc-url"`
	AssetID         string `json:"asset-id"            mapstructure:"asset-id"`
	ShamirThreshold int    `json:"shamir-threshold"    mapstructure:"shamir-threshold"`
	ShamirShares    int    `json:"shamir-shares"       mapstructure:"shamir-shares"`
	VirtualEnvPath  string `json:"virtual-env-path"    mapstructure:"virtual-env-path"`
	TestMode        bool   `json:"test-mode"           mapstructure:"test-mode"`
}

// global singleton
var (
	config     *Config
	initConfig sync.Once
)

// DefaultConfig returns TA's default configuration.
func DefaultConfig() *Config {
	return &Config{
		ServiceBind:     "localhost",
		ServicePort:     8080,
		ShareHolderList: "https://localhost:8081,https://localhost:8082,https://localhost:8083",
		CertsPath:       "./certs/",
		WalletName:      "wallet",
		ElementsRpcUrl:  "http://127.0.0.1:26657",
		AssetID:         "RDDL",
		ShamirThreshold: 2,
		ShamirShares:    3,
		VirtualEnvPath:  "/home/jeckel/develop/rddl/python-shamir-mnemonic/.venv/",
		TestMode:        false,
	}
}

// GetConfig returns the config instance for the SDK.
func GetConfig() *Config {
	initConfig.Do(func() {
		config = DefaultConfig()
	})
	return config
}
