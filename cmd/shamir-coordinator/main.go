package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/rddl-network/shamir-coordinator-service/config"
	"github.com/rddl-network/shamir-coordinator-service/service"
	"github.com/rddl-network/shamir-shareholder-service/client"
)

func main() {
	cfg, err := config.LoadConfig("./")
	if err != nil {
		log.Fatalf("fatal error reading the configuration %s", err)
	}

	// initializing all shareholder clients
	var sscs []client.IShamirShareholderClient
	for _, host := range strings.Split(cfg.ShareHolderList, ",") {
		mTLSClient, err := service.Get2wayTLSClient(cfg)
		if err != nil {
			log.Fatalf("fatal error setting up mutual TLS shareholder client")
		}
		ssc := client.NewShamirShareholderClient(host, mTLSClient)
		sscs = append(sscs, ssc)
	}

	slip39Interface := &service.Slip39Interface{}
	SCoordinator := service.NewShamirCoordinatorService(cfg, sscs, slip39Interface)
	err = SCoordinator.Run()
	if err != nil {
		fmt.Print(err.Error())
	}
}
