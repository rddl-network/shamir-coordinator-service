package service

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/rddl-network/shamir-coordinator-service/config"
)

type ShamirCoordinatorService struct {
	cfg    *config.Config
	router *gin.Engine
}

func NewShamirCoordinatorService(cfg *config.Config) *ShamirCoordinatorService {
	service := &ShamirCoordinatorService{}
	service.cfg = cfg

	gin.SetMode(gin.ReleaseMode)
	service.router = gin.New()
	service.router.POST("/send/:recipient/:amount", service.sendTokens)

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
