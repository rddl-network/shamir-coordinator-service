package service_test

import (
	"testing"

	elements "github.com/rddl-network/elements-rpc"
	elementsmocks "github.com/rddl-network/elements-rpc/utils/mocks"
	"github.com/rddl-network/shamir-coordinator-service/config"
	"github.com/rddl-network/shamir-coordinator-service/service"
	"github.com/rddl-network/shamir-coordinator-service/testutil"
	"github.com/stretchr/testify/assert"
)

func TestTestMode(t *testing.T) {
	cfg, err := config.LoadConfig("../")
	assert.NoError(t, err)
	mycfg := *cfg
	mycfg.TestMode = true

	ssc := testutil.NewShamirShareholderClientMock(&mycfg)
	s := service.NewShamirCoordinatorService(&mycfg, ssc)

	routes := s.GetRoutes()
	assert.Equal(t, 3, len(routes))
}

func TestNotTestMode(t *testing.T) {
	cfg, err := config.LoadConfig("../")
	assert.NoError(t, err)
	ssc := testutil.NewShamirShareholderClientMock(cfg)
	s := service.NewShamirCoordinatorService(cfg, ssc)

	routes := s.GetRoutes()
	assert.Equal(t, 2, len(routes))
}

func TestSendTo(t *testing.T) {
	elements.Client = &elementsmocks.MockClient{}

	cfg, err := config.LoadConfig("../")
	assert.NoError(t, err)
	ssc := testutil.NewShamirShareholderClientMock(cfg)
	s := service.NewShamirCoordinatorService(cfg, ssc)
	address := "tlq1qqvsmfp0w3dmvwtkfteanzk0n7wksu6zx4pywzvak9p6d34yngghw39ynqwcxqrq3muqxffflmprr9exn8ldm79mlkz7dmpy0e"
	amount := "0.0001"
	txID, err := s.SendAsset(address, amount)
	assert.NoError(t, err)
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000000", txID)
}
