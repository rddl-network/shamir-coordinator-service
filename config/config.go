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
rpc-wallet-name="{{ .RPCWalletName }}"
rpc-host="{{ .RPCHost }}"
rpc-port={{ .RPCPort }}
rpc-user="{{ .RPCUser }}"
rpc-password="{{ .RPCPassword }}"
rpc-scheme="{{ .RPCScheme }}"
rpc-enc-timeout={{ .RPCEncTimeout }}
asset-id="{{ .AssetID }}"
shamir-threshold={{ .ShamirThreshold }}
shamir-shares={{ .ShamirShares }}
test-mode={{ .TestMode }}
log-level="{{ .LogLevel }}"
db-path="{{ .DBPath }}"
`

// Config defines TA's top level configuration
type Config struct {
	ServiceBind     string `json:"service-bind"      mapstructure:"service-bind"`
	ServicePort     int    `json:"service-port"      mapstructure:"service-port"`
	ShareHolderList string `json:"share-holder-list" mapstructure:"share-holder-list"`
	CertsPath       string `json:"certs-path"        mapstructure:"certs-path"`
	RPCWalletName   string `json:"rpc-wallet-name"   mapstructure:"rpc-wallet-name"`
	RPCHost         string `json:"rpc-host"          mapstructure:"rpc-host"`
	RPCPort         int    `json:"rpc-port"          mapstructure:"pc-port"`
	RPCUser         string `json:"rpc-user"          mapstructure:"rpc-user"`
	RPCPassword     string `json:"rpc-password"      mapstructure:"rpc-password"`
	RPCScheme       string `json:"rpc-scheme"        mapstructure:"rpc-scheme"`
	RPCEncTimeout   int    `json:"rpc-enc-timeout"   mapstructure:"rpc-enc-timeout"`
	AssetID         string `json:"asset-id"          mapstructure:"asset-id"`
	ShamirThreshold int    `json:"shamir-threshold"  mapstructure:"shamir-threshold"`
	ShamirShares    int    `json:"shamir-shares"     mapstructure:"shamir-shares"`
	TestMode        bool   `json:"test-mode"         mapstructure:"test-mode"`
	LogLevel        string `json:"log-level"         mapstructure:"log-level"`
	DBPath          string `json:"db-path" 			 mapstructure:"db-path"`
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
		RPCWalletName:   "wallet",
		RPCHost:         "localhost",
		RPCPort:         18884,
		RPCUser:         "user",
		RPCPassword:     "password",
		RPCScheme:       "http",
		RPCEncTimeout:   20,
		AssetID:         "asset-id",
		ShamirThreshold: 2,
		ShamirShares:    3,
		TestMode:        false,
		LogLevel:        "info",
		DBPath:          "./data/",
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
	url := c.RPCScheme + "://" + c.RPCUser + ":" + c.RPCPassword
	url = url + "@" + c.RPCHost + ":" + strconv.Itoa(c.RPCPort)
	url = url + "/wallet/" + c.RPCWalletName
	return url
}
