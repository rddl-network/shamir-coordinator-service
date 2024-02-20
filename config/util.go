package config

import (
	"bytes"
	"log"
	"os"
	"text/template"

	"github.com/spf13/viper"
)

func LoadConfig(path string) (cfg *Config, err error) {
	v := viper.New()
	v.AddConfigPath(path)
	v.SetConfigName("app")
	v.SetConfigType("toml")
	v.AutomaticEnv()

	err = v.ReadInConfig()
	if err == nil {
		cfg = GetConfig()
		cfg.ServiceBind = v.GetString("service-bind")
		cfg.ServicePort = v.GetInt("service-port")
		cfg.CertsPath = v.GetString("certs-path")
		cfg.RPCScheme = v.GetString("rpc-scheme")
		cfg.RPCHost = v.GetString("rpc-host")
		cfg.RPCPort = v.GetInt("rpc-port")
		cfg.RPCUser = v.GetString("rpc-user")
		cfg.RPCPassword = v.GetString("rpc-password")
		cfg.RPCWalletName = v.GetString("rpc-wallet-name")
		cfg.VirtualEnvPath = v.GetString("virtual-env-path")
		cfg.ShamirShares = v.GetInt("shamir-shares")
		cfg.ShamirThreshold = v.GetInt("shamir-threshold")
		cfg.AssetID = v.GetString("asset-id")
		cfg.ShareHolderList = v.GetString("share-holder-list")

		if err := viper.Unmarshal(&cfg); err != nil {
			log.Fatal(err)
		}
		return
	}
	log.Println("no config file found")

	tmpl := template.New("appConfigFileTemplate")
	configTemplate, err := tmpl.Parse(DefaultConfigTemplate)
	if err != nil {
		return
	}

	var buffer bytes.Buffer
	if err = configTemplate.Execute(&buffer, GetConfig()); err != nil {
		return
	}

	if err = v.ReadConfig(&buffer); err != nil {
		return
	}
	if err = v.SafeWriteConfig(); err != nil {
		return
	}

	log.Println("default config file created. please adapt it and restart the application. exiting...")
	os.Exit(0)
	return
}
