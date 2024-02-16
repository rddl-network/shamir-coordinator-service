package service_test

import (
	"testing"

	"github.com/rddl-network/shamir-coordinator-service/config"
	"github.com/rddl-network/shamir-coordinator-service/service"
	"gotest.tools/assert"
)

func TestTestnetModeTrue(t *testing.T) {
	cfg := config.DefaultConfig()
	s := service.NewShamirCoordinatorService(cfg)

	routes := s.GetRoutes()
	assert.Equal(t, 2, len(routes))
	assert.Equal(t, "/register/:pubkey", routes[1].Path)
}

func TestTestnetModeFalse(t *testing.T) {
	cfg := config.DefaultConfig()
	s := service.NewShamirCoordinatorService(cfg)

	routes := s.GetRoutes()
	assert.Equal(t, 1, len(routes))
}
