package service

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/rddl-network/shamir-coordinator-service/config"
)

type ShamirCoordinatorService struct {
	cfg    *config.Config
	Router *gin.Engine
	ssc    IShamirShareholderClient
}

func NewShamirCoordinatorService(cfg *config.Config, ssc IShamirShareholderClient) *ShamirCoordinatorService {
	service := &ShamirCoordinatorService{}
	service.cfg = cfg
	service.ssc = ssc

	gin.SetMode(gin.ReleaseMode)
	service.Router = gin.New()
	service.Router.POST("/send/:recipient/:amount", service.SendTokens)
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
