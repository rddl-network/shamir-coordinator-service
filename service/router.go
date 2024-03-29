package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *ShamirCoordinatorService) SendTokens(c *gin.Context) {
	var request SendTokensRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mnemonics, err := s.CollectMnemonics()
	// This code snippet is handling an error scenario in the `sendTokens` function of the
	// `ShamirCoordinatorService`.
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error collecting the shares"})
		return
	}
	passphrase, err := s.RecoverSeed(mnemonics[:s.cfg.ShamirThreshold])
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error computing the seeds: " + err.Error()})
		return
	}

	// prepare the wallet, loading and unlocking
	err = s.PrepareWallet(passphrase)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error loading the wallet " + err.Error()})
		return
	}
	// send asset
	txID, err := s.SendAsset(request.Recipient, request.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error sending/broadcasting the transaction"})
		return
	}

	var resBody SendTokensResponse
	resBody.TxID = txID
	c.JSON(http.StatusOK, resBody)
}

func (s *ShamirCoordinatorService) DeployShares(c *gin.Context) {
	secret := c.Param("secret")
	if !IsValidHex(secret) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "the secret has to be send in valid hex string format"})
		return
	}
	if len(secret) != 32 && len(secret) != 64 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "the secret has to be of length 32 or 64 (16 or 32 byte)"})
		return
	}

	mnemonics, err := s.CreateMnemonics(secret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "share creation failed: " + err.Error()})
		return
	}
	err = s.deployMnemonics(mnemonics)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error sending/broadcasting the transaction"})
		return
	}

	c.JSON(http.StatusOK, "{}")
}

func (s *ShamirCoordinatorService) CollectShares(c *gin.Context) {
	mnemonics, err := s.CollectMnemonics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error collecting the shares"})
		return
	}
	seed, err := s.RecoverSeed(mnemonics[:s.cfg.ShamirThreshold])
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error computing the seeds: " + err.Error()})
		return
	}
	var resBody MnemonicsResponse
	resBody.Mnemonics = mnemonics
	resBody.Seed = seed
	c.JSON(http.StatusOK, resBody)
}

func (s *ShamirCoordinatorService) GetRoutes() gin.RoutesInfo {
	routes := s.Router.Routes()
	return routes
}
