package service

import (
	"fmt"
	"strconv"

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

func (s *ShamirCoordinatorService) LoadWallet() error {
	_, err := elements.LoadWallet(s.cfg.GetRPCConnectionString(), []string{s.cfg.RpcWalletName})
	if err != nil {
		fmt.Println("Error loading the wallet: " + err.Error())
	}
	return err
}

func (s *ShamirCoordinatorService) setWalletPassphrase(passphrase string, timeout int) error {
	err := elements.Walletpassphrase(s.cfg.GetRPCConnectionString(), []string{passphrase, strconv.Itoa(timeout)})
	if err != nil {
		fmt.Println("Error loading the wallet: " + err.Error())
	}
	return err
}

func (s *ShamirCoordinatorService) UnloadWallet() error {
	_, err := elements.UnloadWallet(s.cfg.GetRPCConnectionString(), []string{s.cfg.RpcWalletName})
	if err != nil {
		fmt.Println("Error unloading the wallet: " + err.Error())
	}
	return err

}
