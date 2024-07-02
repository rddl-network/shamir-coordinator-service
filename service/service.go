package service

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/rddl-network/go-utils/logger"
	"github.com/rddl-network/go-utils/tls"
	"github.com/rddl-network/shamir-coordinator-service/config"
	"github.com/rddl-network/shamir-shareholder-service/client"
)

type ShamirCoordinatorService struct {
	cfg             *config.Config
	Router          *gin.Engine
	sscs            map[string]client.IShamirShareholderClient
	slip39Interface ISlip39
	logger          log.AppLogger
}

func NewShamirCoordinatorService(cfg *config.Config, sscs map[string]client.IShamirShareholderClient, slip39Interface ISlip39, logger log.AppLogger) *ShamirCoordinatorService {
	service := &ShamirCoordinatorService{
		cfg:             cfg,
		sscs:            sscs,
		slip39Interface: slip39Interface,
		logger:          logger,
	}

	gin.SetMode(gin.ReleaseMode)
	service.Router = gin.New()
	service.Router.POST("/send", service.SendTokens)
	service.Router.POST("/reissue", service.ReIssue)
	service.Router.POST("/issue-machine-nft", service.IssueMachineNFT)
	service.Router.POST("/mnemonics/:secret", service.DeployShares)
	if cfg.TestMode {
		service.Router.GET("/mnemonics", service.CollectShares)
	}
	return service
}

func (s *ShamirCoordinatorService) Run() (err error) {
	cfg := config.GetConfig()
	caCertFile, err := os.ReadFile(cfg.CertsPath + "ca.crt")
	if err != nil {
		return err
	}

	tlsConfig := tls.Get2WayTLSServer(caCertFile)
	server := &http.Server{
		Addr:      fmt.Sprintf("%s:%d", cfg.ServiceBind, cfg.ServicePort),
		TLSConfig: tlsConfig,
		Handler:   s.Router,
	}

	// workaround to listen on tcp4 and not tcp6
	// https://stackoverflow.com/a/38592286
	ln, err := net.Listen("tcp4", server.Addr)
	if err != nil {
		return err
	}
	defer ln.Close()

	return server.ServeTLS(ln, cfg.CertsPath+"server.crt", cfg.CertsPath+"server.key")
}
