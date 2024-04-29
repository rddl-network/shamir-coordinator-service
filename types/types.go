package types

type SendTokensRequest struct {
	Recipient string `binding:"required" json:"recipient"`
	Amount    string `binding:"required" json:"amount"`
}

type SendTokensResponse struct {
	TxID string `binding:"required" json:"tx-id"`
}

type MnemonicsResponse struct {
	Mnemonics []string `binding:"required" json:"mnemonics"`
	Seed      string   `binding:"required" json:"seed"`
}
