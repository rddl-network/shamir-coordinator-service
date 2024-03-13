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

func TestSendTo(t *testing.T) {
	elements.Client = &elementsmocks.MockClient{}

	cfg, err := config.LoadConfig("../")
	assert.NoError(t, err)
	ssc := testutil.NewShamirShareholderClientMock(cfg)
	slip39mock := &testutil.Slip39Mock{}
	s := service.NewShamirCoordinatorService(cfg, ssc, slip39mock)
	address := "tlq1qqvsmfp0w3dmvwtkfteanzk0n7wksu6zx4pywzvak9p6d34yngghw39ynqwcxqrq3muqxffflmprr9exn8ldm79mlkz7dmpy0e"
	amount := "0.0001"
	txID, err := s.SendAsset(address, amount)
	assert.NoError(t, err)
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000000", txID)
}

func TestPrepareWallet(t *testing.T) {
	elements.Client = &elementsmocks.MockClient{}

	cfg, err := config.LoadConfig("../")
	assert.NoError(t, err)

	ssc := testutil.NewShamirShareholderClientMock(cfg)
	slip39mock := &testutil.Slip39Mock{}
	s := service.NewShamirCoordinatorService(cfg, ssc, slip39mock)

	passphrase := "password"

	err = s.PrepareWallet(passphrase)
	assert.NoError(t, err)
}
