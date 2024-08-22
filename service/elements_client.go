package service

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"

	elements "github.com/rddl-network/elements-rpc"
	strutil "github.com/rddl-network/go-utils/str"
	"github.com/rddl-network/shamir-coordinator-service/types"
)

var (
	// this mutex has to protect all signing and crafting of transactions and their inputs
	// so that UTXOs are not spend twice by accident
	elementsSyncAccess sync.Mutex
)

func isValidAmount(amount string) (valid bool) {
	amount = strings.TrimSpace(amount)
	f, err := strconv.ParseFloat(amount, 64)
	if err == nil && f > 0.0 {
		valid = true
	}
	return
}

func (s *ShamirCoordinatorService) SendAsset(address string, amount string, asset string) (txID string, err error) {
	if asset == "" {
		asset = s.cfg.AssetID
	}
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
		`"` + asset + `"`,
	})
	return
}

func (s *ShamirCoordinatorService) ReissueAsset(asset string, amount string) (txID string, err error) {
	res, err := elements.ReissueAsset(s.cfg.GetRPCConnectionString(), []string{
		`"` + asset + `"`,
		amount,
	})
	if err != nil {
		return
	}
	txID = res.TxID
	return
}

func (s *ShamirCoordinatorService) IsWalletLoaded(rpcURL string, walletname string) (loaded bool, err error) {
	wallets, err := elements.ListWallets(rpcURL, []string{})
	if err != nil {
		return
	}

	loaded = strutil.ContainsString(wallets, walletname)
	return
}

func (s *ShamirCoordinatorService) PrepareWallet(passphrase string) (err error) {
	// the wallet is expected to be loaded, verify if it's loaded
	loaded, err := s.IsWalletLoaded(s.cfg.GetRPCConnectionString(), s.cfg.RPCWalletName)
	if err != nil {
		s.logger.Error("error", "Error listing the wallets: "+err.Error())
	}
	if !loaded {
		// loaded wallet via RPC if not loaded
		_, err = elements.LoadWallet(s.cfg.GetRPCConnectionString(), []string{`"` + s.cfg.RPCWalletName + `"`})
		if err != nil {
			s.logger.Error("error", "Error loading the wallet: "+err.Error())
			return
		}
	}

	_, err = elements.Walletpassphrase(s.cfg.GetRPCConnectionString(), []string{`"` + passphrase + `"`, strconv.Itoa(s.cfg.RPCEncTimeout)})
	if err != nil {
		s.logger.Error("error", "Error decrypting the wallet: "+err.Error())
		return
	}
	return
}

func (s *ShamirCoordinatorService) IssueNFTAsset(name string, machineAddress string, domain string) (assetID string, contract string, hexTx string, err error) {
	url := s.cfg.GetRPCConnectionString()

	address, err := elements.GetNewAddress(url, []string{``})
	if err != nil {
		return
	}

	addressInfo, err := elements.GetAddressInfo(url, []string{`"` + address + `"`})
	if err != nil {
		return
	}

	elementsSyncAccess.Lock()
	defer elementsSyncAccess.Unlock()
	hex, err := elements.CreateRawTransaction(url, []string{`[]`, `[{"data":"00"}]`})
	if err != nil {
		return
	}

	fundRawTransactionResult, err := elements.FundRawTransaction(url, []string{`"` + hex + `"`, `{"feeRate":0.00001000}`})
	if err != nil {
		return
	}

	c := types.Contract{
		Entity: types.Entity{
			Domain: domain,
		},
		IssuerPubkey: addressInfo.Pubkey,
		MachineAddr:  machineAddress,
		Name:         name,
		Precision:    0,
		Version:      0,
	}
	contractBytes, err := json.Marshal(c)
	if err != nil {
		return
	}
	// e.g. {"entity":{"domain":"testnet-assets.rddl.io"}, "issuer_pubkey":"02...}
	contract = string(contractBytes)

	h := sha256.New()
	_, err = h.Write(contractBytes)
	if err != nil {
		return
	}
	// e.g. 7ca8bb403ee5dccddef7b89b163048cf39439553f0402351217a4a03d2224df8
	hash := h.Sum(nil)

	// Reverse hash, e.g. f84d22d2034a7a21512340f053954339cf4830169bb8f7decddce53e40bba87c
	for i, j := 0, len(hash)-1; i < j; i, j = i+1, j-1 {
		hash[i], hash[j] = hash[j], hash[i]
	}

	rawIssueAssetResults, err := elements.RawIssueAsset(url, []string{`"` + fundRawTransactionResult.Hex + `"`,
		`[{"asset_amount":0.00000001, "asset_address":"` + address + `", "blind":false, "contract_hash":"` + fmt.Sprintf("%+x", hash) + `"}]`,
	})
	if err != nil {
		return
	}

	rawIssueAssetResult := rawIssueAssetResults[len(rawIssueAssetResults)-1]
	hex, err = elements.BlindRawTransaction(url, []string{`"` + rawIssueAssetResult.Hex + `"`, `true`, `[]`, `false`})
	if err != nil {
		return
	}
	assetID = rawIssueAssetResult.Asset

	signRawTransactionWithWalletResult, err := elements.SignRawTransactionWithWallet(url, []string{`"` + hex + `"`})
	if err != nil {
		return
	}

	testMempoolAcceptResults, err := elements.TestMempoolAccept(url, []string{`["` + signRawTransactionWithWalletResult.Hex + `"]`})
	if err != nil {
		return
	}

	testMempoolAcceptResult := testMempoolAcceptResults[len(testMempoolAcceptResults)-1]
	if !testMempoolAcceptResult.Allowed {
		err = fmt.Errorf("not accepted by mempool: %+v %+v", testMempoolAcceptResult, signRawTransactionWithWalletResult)
		return
	}

	hex, err = elements.SendRawTransaction(url, []string{`"` + signRawTransactionWithWalletResult.Hex + `"`})
	if err != nil {
		return
	}

	return assetID, contract, hex, err
}

func (s *ShamirCoordinatorService) WalletLock() (err error) {
	url := s.cfg.GetRPCConnectionString()
	return elements.WalletLock(url)
}
