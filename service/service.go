package service

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/rddl-network/shamir-coordinator-service/config"
)

type ShamirCoordinatorService struct {
	cfg    *config.Config
	router *gin.Engine
	ssc    IShamirShareholderClient
}

func NewShamirCoordinatorService(cfg *config.Config, ssc IShamirShareholderClient) *ShamirCoordinatorService {
	service := &ShamirCoordinatorService{}
	service.cfg = cfg
	service.ssc = ssc

	gin.SetMode(gin.ReleaseMode)
	service.router = gin.New()
	service.router.POST("/send/:recipient/:amount", service.sendTokens)
	service.router.POST("/mnemonics/:secret", service.deployShares)
	if cfg.TestMode {
		service.router.GET("/mnemonics", service.collectShares)
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
	err := s.router.Run(addr)

	return err
}
