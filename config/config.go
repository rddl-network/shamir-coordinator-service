package config

import (
	"strconv"
	"sync"
)

const DefaultConfigTemplate = `
SERVICE_BIND="{{ .ServiceBind }}"
SERVICE_PORT={{ .ServicePort }}
SHARE_HOLDER_LIST="{{ .ShareHolderList }}"
CERTS_PATH="{{ .CertsPath }}"
RPC_WALLET_NAME="{{ .RpcWalletName }}"
RPC_HOST="{{ .RpcHost }}"
RPC_PORT={{ .RpcPort }}
RPC_USER="{{ .RpcUser }}"
RPC_PASSWORD="{{ .RpcPassword }}"
RPC_SCHEME="{{ .RpcScheme }}"
RPC_ENC_TIMEOUT={{ .RpcEncTimeout }}
ASSET_ID="{{ .AssetID }}"
SHAMIR_THRESHOLD={{ .ShamirThreshold }}
SHAMIR_SHARES={{ .ShamirShares }}
VIRTUAL_ENV_PATH="{{ .VirtualEnvPath }}"
TEST_MODE={{ .TestMode }}
`

// Config defines TA's top level configuration
type Config struct {
	ServiceBind     string `json:"service-bind"        mapstructure:"SERVICE_BIND"`
	ServicePort     int    `json:"service-port"        mapstructure:"SERVICE_PORT"`
	ShareHolderList string `json:"share-holder-list"   mapstructure:"SHARE_HOLDER_LIST"`
	CertsPath       string `json:"certs-path"          mapstructure:"CERTS_PATH"`
	RpcWalletName   string `json:"rpc-wallet-name"     mapstructure:"RPC_WALLET_NAME"`
	RpcHost         string `json:"rpc-host"            mapstructure:"RPC_HOST"`
	RpcPort         int    `json:"rpc-port"            mapstructure:"RPC_PORT"`
	RpcUser         string `json:"rpc-user"            mapstructure:"RPC_USER"`
	RpcPassword     string `json:"rpc-password"        mapstructure:"RPC_PASSWORD"`
	RpcScheme       string `json:"rpc-scheme"          mapstructure:"RPC_SCHEME"`
	RpcEncTimeout   int    `json:"rpc-enc-timeout      mapstructure:"RPC_ENC_TIMEOUT"`
	AssetID         string `json:"asset-id"            mapstructure:"ASSET_ID"`
	ShamirThreshold int    `json:"shamir-threshold"    mapstructure:"SHAMIR_THRESHOLD"`
	ShamirShares    int    `json:"shamir-shares"       mapstructure:"SHAMIR_SHARES"`
	VirtualEnvPath  string `json:"virtual-env-path"    mapstructure:"VIRTUAL_ENV_PATH"`
	TestMode        bool   `json:"test-mode"           mapstructure:"TEST_MODE"`
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
		RpcWalletName:   "wallet",
		RpcHost:         "localhost",
		RpcPort:         18884,
		RpcUser:         "user",
		RpcPassword:     "password",
		RpcScheme:       "http",
		RpcEncTimeout:   20,
		AssetID:         "asset-id",
		ShamirThreshold: 2,
		ShamirShares:    3,
		VirtualEnvPath:  "~/.venv/",
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

func (c *Config) GetRPCConnectionString() string {
	url := c.RpcScheme + "://" + c.RpcUser + ":" + c.RpcPassword
	url = url + "@" + c.RpcHost + ":" + strconv.Itoa(c.RpcPort)
	url = url + "/wallet/" + c.RpcWalletName
	return url
}
