package testutil

import (
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rddl-network/go-logger"
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

	return service.NewShamirCoordinatorService(cfg, sscs, slip39Mock, logger.GetLogger(logger.DEBUG))
}

func SetupTestServiceWithSlip39Interface(t *testing.T) *service.ShamirCoordinatorService {
	cfg, err := config.LoadConfig("../")
	assert.NoError(t, err)

	sscs := createShamirShareholderMocks(t, cfg.ShamirShares)

	slip39Mock := &service.Slip39Interface{}

	return service.NewShamirCoordinatorService(cfg, sscs, slip39Mock, logger.GetLogger(logger.DEBUG))
}

func createShamirShareholderMocks(t *testing.T, n int) map[string]client.IShamirShareholderClient {
	ctrl := gomock.NewController(t)
	sscs := make(map[string]client.IShamirShareholderClient)

	for i := 0; i < n; i++ {
		ssc := NewMockIShamirShareholderClient(ctrl)
		ssc.EXPECT().GetMnemonic(gomock.Any()).Return(shareholder.MnemonicBody{}, nil).AnyTimes()
		ssc.EXPECT().PostMnemonic(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
		sscs[strconv.Itoa(i)] = ssc
	}

	return sscs
}
