package service

import (
	elements "github.com/rddl-network/elements-rpc"
)

func (s *ShamirCoordinatorService) SendAsset(address string, amount string) (txID string, err error) {
	txID, err = elements.SendToAddress(s.cfg.GetRPCConnectionString(), []string{
		address,
		`"` + amount + `"`,
		`""`,
		`""`,
		"false",
		"true",
		"null",
		`"unset"`,
		"false",
		`"` + s.cfg.AssetID + `"`,
	})
	return
}
