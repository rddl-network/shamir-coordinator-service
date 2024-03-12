package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type TxIDBody struct {
	TxID string `binding:"required" json:"tx-id"`
}

func (s *ShamirCoordinatorService) sendTokens(c *gin.Context) {
	recipient := c.Param("recipient")
	amount := c.Param("amount")

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
	txID, err := s.SendAsset(recipient, amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error sending/broadcasting the transaction"})
		return
	}

	var resBody TxIDBody
	resBody.TxID = txID
	c.JSON(http.StatusOK, resBody)
}

func (s *ShamirCoordinatorService) deployShares(c *gin.Context) {
	secret := c.Param("secret")
	if !IsValidHex(secret) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "the secret has to be send in valid hex string format"})
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

type MnemonicsBody struct {
	Mnemonics []string `binding:"required" json:"mnemonics"`
	Seed      string   `binding:"required" json:"seed"`
}

func (s *ShamirCoordinatorService) collectShares(c *gin.Context) {
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
	var resBody MnemonicsBody
	resBody.Mnemonics = mnemonics
	resBody.Seed = seed
	c.JSON(http.StatusOK, resBody)
}

func (s *ShamirCoordinatorService) GetRoutes() gin.RoutesInfo {
	routes := s.router.Routes()
	return routes
}
