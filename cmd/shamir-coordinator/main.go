package main

import (
	"fmt"
	"log"

	"github.com/rddl-network/shamir-coordinator-service/config"
	"github.com/rddl-network/shamir-coordinator-service/service"
)

func main() {
	cfg, err := config.LoadConfig("./")
	if err != nil {
		log.Fatalf("fatal error reading the configuration %s", err)
	}
	ssc := service.NewShamirShareholderClient(cfg)
	SCoordinator := service.NewShamirCoordinatorService(cfg, ssc)
	err = SCoordinator.Run()
	if err != nil {
		fmt.Print(err.Error())
	}
}
