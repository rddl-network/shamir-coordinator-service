package service

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	elements "github.com/rddl-network/elements-rpc"
)

type TxIDBody struct {
	TxID string `binding:"required" json:"txID"`
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

	// load wallet
	// decrypt loaded wallet via RPC and above recovered key
	_, err = elements.LoadWallet(s.cfg.GetRPCConnectionString(), []string{s.cfg.RpcWalletName})
	if err != nil {
		fmt.Println("Error loading the wallet: " + err.Error())
		return
	}

	err = elements.Walletpassphrase(s.cfg.GetRPCConnectionString(), []string{passphrase, strconv.Itoa(s.cfg.RpcEncTimeout)})
	if err != nil {
		fmt.Println("Error decrypting the wallet: " + err.Error())
		return
	}

	txID, err := s.SendAsset(recipient, amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error sending/broadcasting the transaction"})
		return
	}

	// unload wallet
	_, err = elements.UnloadWallet(s.cfg.GetRPCConnectionString(), []string{s.cfg.RpcWalletName})
	if err != nil {
		fmt.Println("Error unloading the wallet: " + err.Error())
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
