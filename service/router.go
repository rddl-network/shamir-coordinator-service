package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	hexutil "github.com/rddl-network/go-utils/hex"
	"github.com/rddl-network/shamir-coordinator-service/types"
)

const (
	errCompMsg = "error computing the seeds: "
)

func (s *ShamirCoordinatorService) SendTokens(c *gin.Context) {
	var request types.SendTokensRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}
	s.logger.Info("msg", "preparing to send "+request.Amount+" tokens to "+request.Recipient)
	mnemonics, err := s.CollectMnemonics()
	// This code snippet is handling an error scenario in the `sendTokens` function of the
	// `ShamirCoordinatorService`.
	if err != nil {
		s.logger.Error("error", "error collecting the shares: "+err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "error collecting the shares"})
		return
	}
	passphrase, err := s.RecoverSeed(mnemonics[:s.cfg.ShamirThreshold])
	if err != nil {
		s.logger.Error("error", errCompMsg+err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"Error": errCompMsg + err.Error()})
		return
	}

	// prepare the wallet, loading and unlocking
	err = s.PrepareWallet(passphrase)
	if err != nil {
		s.logger.Error("error", "error loading the wallet: "+err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "error loading the wallet " + err.Error()})
		return
	}
	// send asset
	txID, err := s.SendAsset(request.Recipient, request.Amount)
	if err != nil {
		s.logger.Error("error", "error sending the transaction: "+err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "error sending/broadcasting the transaction"})
		return
	}

	s.logger.Info("msg", "successfully sended tx with id : "+txID+" to "+request.Recipient)
	var resBody types.SendTokensResponse
	resBody.TxID = txID
	c.JSON(http.StatusOK, resBody)
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
