package testutil

import (
	"github.com/rddl-network/shamir-coordinator-service/config"
)

type ShamirShareholderClientMock struct {
	cfg *config.Config
}

func NewShamirShareholderClientMock(cfg *config.Config) *ShamirShareholderClientMock {
	ssc := &ShamirShareholderClientMock{}
	ssc.cfg = cfg
	return ssc
}

func (s *ShamirShareholderClientMock) GetMnemonic(shareholderURI string) (mnemonic string, err error) {
	return
}

func (s *ShamirShareholderClientMock) PostMnemonic(shareHolderUri string, mnemonic string) (err error) {
	return
}
