package service

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	hexutil "github.com/rddl-network/go-utils/hex"
	"github.com/rddl-network/shamir-coordinator-service/types"
)

const (
	errCompMsg       = "error computing the seeds: "
	errWalletMsg     = "error loading the wallet: "
	errSendingTxMsg  = "error sending the transaction: "
	errWalletLockMsg = "error locking wallet: "
)

func (s *ShamirCoordinatorService) AddToQueue(err error) bool {
	// Invalid Bitcoin address response error
	if strings.Contains(err.Error(), "Invalid Bitcoin address:") || !strings.HasSuffix(err.Error(), ": -5") {
		return false
	}
	return true
}

func (s *ShamirCoordinatorService) SendTokens(c *gin.Context) {
	var request types.SendTokensRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}
	s.logger.Info("msg", "preparing to send "+request.Amount+" tokens to "+request.Recipient)
	passphrase, err := s.GetPassphrase()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	// prepare the wallet, loading and unlocking
	err = s.PrepareWallet(passphrase)
	if err != nil {
		s.logger.Error("error", errWalletMsg+err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"Error": errWalletMsg + err.Error()})
		return
	}
	// send asset
	txID, err := s.SendAsset(request.Recipient, request.Amount, request.Asset)
	if err != nil {
		s.logger.Error("error", errSendingTxMsg+err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "error sending/broadcasting the transaction"})
		if s.AddToQueue(err) {
			if e := s.db.CreateSendTokensRequest(request.Recipient, request.Amount, request.Asset); e != nil {
				s.logger.Error("error", "error storing transaction request: "+e.Error())
			}
		}
	} else {
		s.logger.Info("msg", "successfully sended tx with id: "+txID+" to "+request.Recipient)
		var resBody types.SendTokensResponse
		resBody.TxID = txID
		c.JSON(http.StatusOK, resBody)
	}

	if err = s.WalletLock(); err != nil {
		s.logger.Error("error", errWalletLockMsg+err.Error())
	}
}

func (s *ShamirCoordinatorService) ReIssue(c *gin.Context) {
	var request types.ReIssueRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}
	s.logger.Info("msg", "preparing to reissue "+request.Amount+" of asset "+request.Asset)

	passphrase, err := s.GetPassphrase()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	// prepare the wallet, loading and unlocking
	err = s.PrepareWallet(passphrase)
	if err != nil {
		s.logger.Error("error", errWalletMsg+err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"Error": errWalletMsg + err.Error()})
		return
	}

	// reissue asset
	txID, err := s.ReissueAsset(request.Asset, request.Amount)
	if err != nil {
		s.logger.Error("error", "error reissuing asset: "+err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "error reissuing asset"})
		if e := s.db.CreateReIssueRequest(request.Amount, request.Asset); e != nil {
			s.logger.Error("error", "error storing reissue request: "+e.Error())
		}
	} else {
		s.logger.Info("msg", "successfully reissued asset", "tx-id", txID, "asset", request.Asset, "amount", request.Amount)
		c.JSON(http.StatusOK, types.ReIssueResponse{TxID: txID})
	}

	if err = s.WalletLock(); err != nil {
		s.logger.Error("error", errWalletLockMsg+err.Error())
	}
}

func (s *ShamirCoordinatorService) IssueMachineNFT(c *gin.Context) {
	var request types.IssueMachineNFTRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	s.logger.Info("msg", "preparing to issue machine nft", "name", request.Name, "machineAddress", request.MachineAddress, "domain", request.Domain)

	passphrase, err := s.GetPassphrase()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	// prepare the wallet, loading and unlocking
	err = s.PrepareWallet(passphrase)
	if err != nil {
		s.logger.Error("error", errWalletMsg+err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"Error": errWalletMsg + err.Error()})
		return
	}

	asset, contract, hexTx, err := s.IssueNFTAsset(request.Name, request.MachineAddress, request.Domain)
	if err != nil {
		s.logger.Error("error", "error issuing machine nft: "+err.Error(), "name", request.Name, "machineAddress", request.MachineAddress, "domain", request.Domain)
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		if e := s.db.CreateIssueMachineNFTRequest(request.Name, request.MachineAddress, request.Domain); e != nil {
			s.logger.Error("error", "error storing issue nft request: "+e.Error())
		}
	} else {
		s.logger.Info("msg", "successfully issued machine nft", "asset_id", asset, "contract", contract, "hex_tx", hexTx)
		c.JSON(http.StatusOK, types.IssueMachineNFTResponse{
			Asset:    asset,
			Contract: contract,
			HexTX:    hexTx,
		})
	}

	if err = s.WalletLock(); err != nil {
		s.logger.Error("error", errWalletLockMsg+err.Error())
	}
}

func (s *ShamirCoordinatorService) DeployShares(c *gin.Context) {
	secret := c.Param("secret")
	if !hexutil.IsValidHex(secret) {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "the secret has to be send in valid hex string format"})
		return
	}
	if len(secret) != 32 && len(secret) != 64 {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "the secret has to be of length 32 or 64 (16 or 32 byte)"})
		return
	}

	mnemonics, err := s.CreateMnemonics(secret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "share creation failed: " + err.Error()})
		return
	}
	err = s.deployMnemonics(mnemonics)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "error sending/broadcasting the transaction"})
		return
	}

	c.JSON(http.StatusOK, "{}")
}

func (s *ShamirCoordinatorService) CollectShares(c *gin.Context) {
	mnemonics, err := s.CollectMnemonics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "error collecting the shares"})
		return
	}
	seed, err := s.RecoverSeed(mnemonics[:s.cfg.ShamirThreshold])
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": errCompMsg + err.Error()})
		return
	}
	var resBody types.MnemonicsResponse
	resBody.Mnemonics = mnemonics
	resBody.Seed = seed
	c.JSON(http.StatusOK, resBody)
}

func (s *ShamirCoordinatorService) GetRoutes() gin.RoutesInfo {
	routes := s.Router.Routes()
	return routes
}
