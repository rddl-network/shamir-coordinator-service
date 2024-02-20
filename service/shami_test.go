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

	seed := "31622fc2d536a751dfff93c6cf21b3d206d4c5362f7fa48e974233db0a56c6c73e6c7466e424d1fd04ed5e0e94e155ad"
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

	var mnemonics = []string{"warmth merchant academic acid admit slice steady index prevent counter unusual fishing fatal depend fawn class square depend holy fatigue mixed salon beard omit cause mixture pistol station object frozen privacy visual theory sister teacher treat platform painting exercise employer emission favorite devote angel voice center",
		"warmth merchant academic agency adjust luck trust angry device writing flavor emperor payment reunion crisis olympic desire treat keyboard ajar actress practice single unhappy lobe robin agency rescue military capacity liquid railroad smart harvest prize random spray domestic hand problem class museum laundry debris withdraw pencil",
		"warmth merchant academic always adult envy disease legend orange literary strike husband switch duke crunch ending kernel coal mayor yoga public provide hazard isolate guest guest island award therapy hand review imply spit leaves kind fake drug lilac loud lunch medal genre perfect beard spray exhaust"}
	seed, err := s.RecoverSeed(mnemonics[:cfg.ShamirThreshold])
	assert.NoError(t, err)
	assert.Equal(t, "31622fc2d536a751dfff93c6cf21b3d206d4c5362f7fa48e974233db0a56c6c73e6c7466e424d1fd04ed5e0e94e155ad", seed)
}
