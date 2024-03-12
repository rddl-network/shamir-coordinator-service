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

func TestPrepareWallet(t *testing.T) {
	elements.Client = &elementsmocks.MockClient{}

	cfg, err := config.LoadConfig("../")
	assert.NoError(t, err)

	ssc := testutil.NewShamirShareholderClientMock(cfg)
	s := service.NewShamirCoordinatorService(cfg, ssc)

	passphrase := "password"

	err = s.PrepareWallet(passphrase)
	assert.NoError(t, err)
}
