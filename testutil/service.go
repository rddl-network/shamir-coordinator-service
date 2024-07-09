package testutil

import (
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	log "github.com/rddl-network/go-utils/logger"
	"github.com/rddl-network/shamir-coordinator-service/config"
	"github.com/rddl-network/shamir-coordinator-service/service"
	"github.com/rddl-network/shamir-coordinator-service/service/backend"
	"github.com/rddl-network/shamir-shareholder-service/client"
	shareholder "github.com/rddl-network/shamir-shareholder-service/service"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/storage"
)

func SetupTestService(t *testing.T) *service.ShamirCoordinatorService {
	cfg := config.GetConfig()
	sscs := createShamirShareholderMocks(t, cfg.ShamirShares)

	slip39Mock := &Slip39Mock{}
	logger := log.GetLogger(log.DEBUG)
	db := SetupTestDBConnector(t)

	return service.NewShamirCoordinatorService(cfg, sscs, slip39Mock, logger, db)
}

func SetupTestServiceWithSlip39Interface(t *testing.T) *service.ShamirCoordinatorService {
	cfg := config.GetConfig()
	sscs := createShamirShareholderMocks(t, cfg.ShamirShares)

	slip39Mock := &service.Slip39Interface{}
	logger := log.GetLogger(log.DEBUG)
	db := SetupTestDBConnector(t)

	return service.NewShamirCoordinatorService(cfg, sscs, slip39Mock, logger, db)
}

func SetupTestDBConnector(t *testing.T) *backend.DBConnector {
	db, err := leveldb.Open(storage.NewMemStorage(), nil)
	if err != nil {
		t.Fatal("Error opening in-memory LevelDB: ", err)
	}
	return backend.NewDBConnector(db)
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
