package service_test

import (
	"testing"

	"github.com/rddl-network/shamir-coordinator-service/config"
	"github.com/rddl-network/shamir-coordinator-service/service"
	"github.com/rddl-network/shamir-coordinator-service/testutil"
	"github.com/stretchr/testify/assert"
)

func TestShamirDeploymnet(t *testing.T) {
	cfg, err := config.LoadConfig("../")
	assert.NoError(t, err)
	ssc := testutil.NewShamirShareholderClientMock(cfg)
	s := service.NewShamirCoordinatorService(cfg, ssc)

	seed := "31622fc2d536a751dfff93c6cf21b3d206d4c5362f7fa48e974233db0a56c6c7"
	mnemonics, err := s.CreateMnemonics(seed)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(mnemonics))

	recoveredSeed, err := s.RecoverSeed(mnemonics[:cfg.ShamirThreshold])
	assert.NoError(t, err)
	assert.Equal(t, seed, recoveredSeed)
}

func TestShamirRecovery(t *testing.T) {
	cfg, err := config.LoadConfig("../")
	assert.NoError(t, err)
	ssc := testutil.NewShamirShareholderClientMock(cfg)
	s := service.NewShamirCoordinatorService(cfg, ssc)

	var mnemonics = []string{"military upgrade academic acid agency grasp superior empty bundle network wrist plot raisin identify ranked install segment email calcium view fragment pitch obtain realize costume emission roster toxic airport imply cleanup canyon grownup",
		"military upgrade academic always aviation listen reunion wireless regret work distance else crazy brother modify union cards crazy crucial story jacket invasion mailman fantasy agree marathon view activity pistol provide snake window romantic",
	}
	seed, err := s.RecoverSeed(mnemonics[:cfg.ShamirThreshold])
	assert.NoError(t, err)
	assert.Equal(t, "31622fc2d536a751dfff93c6cf21b3d206d4c5362f7fa48e974233db0a56c6c7", seed)
}
