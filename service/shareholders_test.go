package service_test

import (
	"testing"

	"github.com/rddl-network/shamir-coordinator-service/config"
	"github.com/rddl-network/shamir-coordinator-service/service"
	"github.com/rddl-network/shamir-coordinator-service/testutil"
	"github.com/stretchr/testify/assert"
)

func TestCollectMnemonics(t *testing.T) {
	cfg := config.DefaultConfig()
	ssc := testutil.NewShamirShareholderClientMock(cfg)
	s := service.NewShamirCoordinatorService(cfg, ssc)

	mnemonics, err := s.CollectMnemonics()
	assert.NoError(t, err)
	assert.Equal(t, 3, len(mnemonics))

}
