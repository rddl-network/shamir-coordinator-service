package config

import "sync"

const DefaultConfigTemplate = `
service-host="{{ .ServiceHost }}"
service-port={{ .ServicePort }}
share-holder-list="{{ .ShareHolderList }}"
certs-path="{{ .CertsPath }}"
wallet-name="{{ .WalletName }}"
elements-rpc-url="{{ .ElementsRpcUrl }}"
token="{{ .Token }}"
`

// Config defines TA's top level configuration
type Config struct {
	ServiceBind     string `json:"service-bind"        mapstructure:"service-bind"`
	ServicePort     int    `json:"service-port"        mapstructure:"service-port"`
	ShareHolderList string `json:"share-holder-list"   mapstructure:"share-holder-list"`
	CertsPath       string `json:"certs-path"          mapstructure:"certs-path"`
	WalletName      string `json:"wallet-name"         mapstructure:"wallet-name"`
	ElementsRpcUrl  string `json:"elements-rpc-url"    mapstructure:"elements-rpc-url"`
	Token           string `json:"token"    		   mapstructure:"token"`
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
		ShareHolderList: "localhost:8080",
		CertsPath:       "./certs/",
		WalletName:      "wallet",
		ElementsRpcUrl:  "http://127.0.0.1:26657",
		Token:           "RDDL",
	}
}

// GetConfig returns the config instance for the SDK.
func GetConfig() *Config {
	initConfig.Do(func() {
		config = DefaultConfig()
	})
	return config
}
