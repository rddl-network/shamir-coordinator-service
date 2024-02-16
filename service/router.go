package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type TxIDBody struct {
	TxID string `binding:"required" json:"txID"`
}

func (s *ShamirCoordinatorService) sendTokens(c *gin.Context) {
	recipient := c.Param("recipient")
	amount := c.Param("amount")

	txID, err := s.SendAsset(recipient, amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error sending/broadcasting the transaction"})
		return
	}

	var resBody TxIDBody
	resBody.TxID = txID
	c.JSON(http.StatusOK, resBody)
}

func (s *ShamirCoordinatorService) GetRoutes() gin.RoutesInfo {
	routes := s.router.Routes()
	return routes
}
