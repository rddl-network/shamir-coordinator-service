package testutil

import (
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	log "github.com/rddl-network/go-utils/logger"
	"github.com/rddl-network/shamir-coordinator-service/config"
	"github.com/rddl-network/shamir-coordinator-service/service"
	"github.com/rddl-network/shamir-shareholder-service/client"
	shareholder "github.com/rddl-network/shamir-shareholder-service/service"
)

func SetupTestService(t *testing.T) *service.ShamirCoordinatorService {
	cfg := config.GetConfig()
	sscs := createShamirShareholderMocks(t, cfg.ShamirShares)

	slip39Mock := &Slip39Mock{}
	logger := log.GetLogger(log.DEBUG)

	return service.NewShamirCoordinatorService(cfg, sscs, slip39Mock, logger)
}

func SetupTestServiceWithSlip39Interface(t *testing.T) *service.ShamirCoordinatorService {
	cfg := config.GetConfig()
	sscs := createShamirShareholderMocks(t, cfg.ShamirShares)

	slip39Mock := &service.Slip39Interface{}
	logger := log.GetLogger(log.DEBUG)

	return service.NewShamirCoordinatorService(cfg, sscs, slip39Mock, logger)
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
