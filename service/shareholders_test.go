package service_test

import (
	"testing"

	"github.com/rddl-network/shamir-coordinator-service/config"
	"github.com/rddl-network/shamir-coordinator-service/service"
	"github.com/stretchr/testify/assert"
)

func TestCollectMnemonics(t *testing.T) {
	cfg := config.DefaultConfig()
	s := service.NewShamirCoordinatorService(cfg)

	mnemonics, err := s.CollectMnemonics()
	assert.NoError(t, err)
	assert.Equal(t, 3, len(mnemonics))

}
