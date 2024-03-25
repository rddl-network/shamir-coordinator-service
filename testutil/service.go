package testutil

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rddl-network/shamir-coordinator-service/config"
	"github.com/rddl-network/shamir-coordinator-service/service"
	"github.com/rddl-network/shamir-shareholder-service/client"
	shareholder "github.com/rddl-network/shamir-shareholder-service/service"
	"github.com/stretchr/testify/assert"
)

func SetupTestService(t *testing.T) *service.ShamirCoordinatorService {
	cfg, err := config.LoadConfig("../")
	assert.NoError(t, err)

	sscs := createShamirShareholderMocks(t, cfg.ShamirShares)

	slip39Mock := &Slip39Mock{}

	return service.NewShamirCoordinatorService(cfg, sscs, slip39Mock)
}

func SetupTestServiceWithSlip39Interface(t *testing.T) *service.ShamirCoordinatorService {
	cfg, err := config.LoadConfig("../")
	assert.NoError(t, err)

	sscs := createShamirShareholderMocks(t, cfg.ShamirShares)

	slip39Mock := &service.Slip39Interface{}

	return service.NewShamirCoordinatorService(cfg, sscs, slip39Mock)
}

func createShamirShareholderMocks(t *testing.T, n int) []client.IShamirShareholderClient {
	ctrl := gomock.NewController(t)
	sscs := make([]client.IShamirShareholderClient, n)

	for i := range sscs {
		ssc := NewMockIShamirShareholderClient(ctrl)
		ssc.EXPECT().GetMnemonic(gomock.Any()).Return(shareholder.MnemonicBody{}, nil).AnyTimes()
		ssc.EXPECT().PostMnemonic(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
		sscs[i] = ssc
	}

	return sscs
}
