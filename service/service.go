package service

import (
	"fmt"

	"github.com/gin-gonic/gin"
	log "github.com/rddl-network/go-utils/logger"
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
	service.Router.POST("/mnemonics/:secret", service.DeployShares)
	if cfg.TestMode {
		service.Router.GET("/mnemonics", service.CollectShares)
	}
	return service
}

func (s *ShamirCoordinatorService) Run() (err error) {
	err = s.startWebService()
	if err != nil {
		s.logger.Error("error", err.Error())
	}
	return err
}

func (s *ShamirCoordinatorService) startWebService() error {
	addr := fmt.Sprintf("%s:%d", s.cfg.ServiceBind, s.cfg.ServicePort)
	err := s.Router.Run(addr)

	return err
}
