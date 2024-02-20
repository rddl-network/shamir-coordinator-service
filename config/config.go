package config

import (
	"strconv"
	"sync"
)

const DefaultConfigTemplate = `
service-bind="{{ .ServiceBind }}"
service-port={{ .ServicePort }}
share-holder-list="{{ .ShareHolderList }}"
certs-path="{{ .CertsPath }}"
rpc-wallet-name="{{ .RpcWalletName }}"
rpc-host="{{ .RpcHost }}"
rpc-port={{ .RpcPort }}
rpc-user="{{ .RpcUser }}"
rpc-password="{{ .RpcPassword }}"
rpc-scheme="{{ .RpcScheme }}"
rpc-enc-timeout={{ .RpcEncTimeout }}
asset-id="{{ .AssetID }}"
shamir-threshold={{ .ShamirThreshold }}
shamir-shares={{ .ShamirShares }}
virtual-env-path="{{ .VirtualEnvPath }}"
test-mode={{ .TestMode }}
`

// Config defines TA's top level configuration
type Config struct {
	ServiceBind     string `json:"service-bind"        mapstructure:"service-bind"`
	ServicePort     int    `json:"service-port"        mapstructure:"service-port"`
	ShareHolderList string `json:"share-holder-list"   mapstructure:"share-holder-list"`
	CertsPath       string `json:"certs-path"          mapstructure:"certs-path"`
	RpcWalletName   string `json:"rpc-wallet-name"     mapstructure:"rpc-wallet-name"`
	RpcHost         string `json:"rpc-host"            mapstructure:"rpc-host"`
	RpcPort         int    `json:"rpc-port"            mapstructure:"pc-port"`
	RpcUser         string `json:"rpc-user"            mapstructure:"rpc-user"`
	RpcPassword     string `json:"rpc-password"        mapstructure:"rpc-password"`
	RpcScheme       string `json:"rpc-scheme"          mapstructure:"rpc-scheme"`
	RpcEncTimeout   int    `json:"rpc-enc-timeout"     mapstructure:"rpc-enc-timeout"`
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
		VirtualEnvPath:  "/opt/hostedtoolcache/Python/3.10.13/x64",
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
