package main

import (
	stdlog "log"
	"strings"

	log "github.com/rddl-network/go-utils/logger"
	"github.com/rddl-network/shamir-coordinator-service/config"
	"github.com/rddl-network/shamir-coordinator-service/service"
	"github.com/rddl-network/shamir-shareholder-service/client"
)

func main() {
	cfg, err := config.LoadConfig("./")
	if err != nil {
		stdlog.Fatalf("fatal error reading the configuration %s", err)
	}

	// initializing all shareholder clients
	shareholderHosts := strings.Split(cfg.ShareHolderList, ",")
	sscs := make(map[string]client.IShamirShareholderClient)
	for _, host := range shareholderHosts {
		mTLSClient, err := service.Get2wayTLSClient(cfg)
		if err != nil {
			stdlog.Fatalf("fatal error setting up mutual TLS shareholder client")
		}
		ssc := client.NewShamirShareholderClient(host, mTLSClient)
		sscs[host] = ssc
	}

	slip39Interface := &service.Slip39Interface{}
	logger := log.GetLogger(cfg.LogLevel)
	SCoordinator := service.NewShamirCoordinatorService(cfg, sscs, slip39Interface, logger)
	err = SCoordinator.Run()
	if err != nil {
		logger.Error(err.Error())
	}
}
