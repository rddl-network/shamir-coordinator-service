package service_test

import (
	"testing"

	"github.com/rddl-network/shamir-coordinator-service/config"
	"github.com/rddl-network/shamir-coordinator-service/service"
	"github.com/rddl-network/shamir-coordinator-service/testutil"
	"github.com/stretchr/testify/assert"
)

func TestCollectMnemonics(t *testing.T) {
	cfg, err := config.LoadConfig("../")
	assert.NoError(t, err)
	ssc := testutil.NewShamirShareholderClientMock(cfg)
	slip39mock := &testutil.Slip39Mock{}
	s := service.NewShamirCoordinatorService(cfg, ssc, slip39mock)

	mnemonics, err := s.CollectMnemonics()
	assert.NoError(t, err)
	assert.Equal(t, 3, len(mnemonics))
}
