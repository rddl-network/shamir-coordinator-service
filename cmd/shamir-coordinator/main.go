package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"

	"github.com/rddl-network/shamir-coordinator-service/config"
	"github.com/rddl-network/shamir-coordinator-service/service"

	"github.com/spf13/viper"
)

func loadConfig(path string) (cfg *config.Config, err error) {
	v := viper.New()
	v.AddConfigPath(path)
	v.SetConfigName("app")
	v.SetConfigType("env")
	v.AutomaticEnv()

	err = v.ReadInConfig()
	if err == nil {
		cfg = config.GetConfig()
		cfg.ServiceBind = v.GetString("SERVICE_BIND")
		cfg.ServicePort = v.GetInt("SERVICE_PORT")
		return
	}
	log.Println("no config file found")

	tmpl := template.New("appConfigFileTemplate")
	configTemplate, err := tmpl.Parse(config.DefaultConfigTemplate)
	if err != nil {
		return
	}

	var buffer bytes.Buffer
	if err = configTemplate.Execute(&buffer, config.GetConfig()); err != nil {
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

func main() {
	cfg, err := loadConfig("./")
	if err != nil {
		log.Fatalf("fatal error reading the configuration %s", err)
	}

	SCoordinator := service.NewShamirCoordinatorService(cfg)
	err = SCoordinator.Run()
	if err != nil {
		fmt.Print(err.Error())
	}
}
