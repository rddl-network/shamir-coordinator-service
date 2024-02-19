package service_test

import (
	"testing"

	"github.com/rddl-network/shamir-coordinator-service/config"
	"github.com/rddl-network/shamir-coordinator-service/service"
	"gotest.tools/assert"
)

func TestTestMode(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.TestMode = true
	s := service.NewShamirCoordinatorService(cfg)

	routes := s.GetRoutes()
	assert.Equal(t, 3, len(routes))
}

func TestNotTestMode(t *testing.T) {
	cfg := config.DefaultConfig()
	s := service.NewShamirCoordinatorService(cfg)

	routes := s.GetRoutes()
	assert.Equal(t, 1, len(routes))
}
