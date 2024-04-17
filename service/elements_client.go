package service

import (
	"strconv"

	elements "github.com/rddl-network/elements-rpc"
)

func (s *ShamirCoordinatorService) SendAsset(address string, amount string) (txID string, err error) {
	txID, err = elements.SendToAddress(s.cfg.GetRPCConnectionString(), []string{
		`"` + address + `"`,
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

func (s *ShamirCoordinatorService) IsWalletLoaded(rpcURL string, walletname string) (loaded bool, err error) {
	wallets, err := elements.ListWallets(rpcURL, []string{})
	if err != nil {
		return
	}

	loaded = ContainsString(wallets, walletname)
	return
}

func (s *ShamirCoordinatorService) PrepareWallet(passphrase string) (err error) {
	// the wallet is expected to be loaded, verify if it's loaded
	loaded, err := s.IsWalletLoaded(s.cfg.GetRPCConnectionString(), s.cfg.RPCWalletName)
	if err != nil {
		s.logger.Error("msg", "Error listing the wallets: "+err.Error())
	}
	if !loaded {
		// loaded wallet via RPC if not loaded
		_, err = elements.LoadWallet(s.cfg.GetRPCConnectionString(), []string{`"` + s.cfg.RPCWalletName + `"`})
		if err != nil {
			s.logger.Error("msg", "Error loading the wallet: "+err.Error())
			return
		}
	}

	_, err = elements.Walletpassphrase(s.cfg.GetRPCConnectionString(), []string{`"` + passphrase + `"`, strconv.Itoa(s.cfg.RPCEncTimeout)})
	if err != nil {
		s.logger.Error("msg", "Error decrypting the wallet: "+err.Error())
		return
	}
	return
}
