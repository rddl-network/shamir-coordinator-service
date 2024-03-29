package service

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/rddl-network/shamir-coordinator-service/config"
	"github.com/rddl-network/shamir-shareholder-service/client"
)

type ShamirCoordinatorService struct {
	cfg             *config.Config
	Router          *gin.Engine
	sscs            map[string]client.IShamirShareholderClient
	slip39Interface ISlip39
}

func NewShamirCoordinatorService(cfg *config.Config, sscs map[string]client.IShamirShareholderClient, slip39Interface ISlip39) *ShamirCoordinatorService {
	service := &ShamirCoordinatorService{}
	service.cfg = cfg
	service.sscs = sscs
	service.slip39Interface = slip39Interface

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
		fmt.Print(err.Error())
	}
	return err
}

func (s *ShamirCoordinatorService) startWebService() error {
	addr := fmt.Sprintf("%s:%d", s.cfg.ServiceBind, s.cfg.ServicePort)
	err := s.Router.Run(addr)

	return err
}
