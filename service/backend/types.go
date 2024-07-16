package backend

type SendTokensRequest struct {
	Recipient string `binding:"required" json:"recipient"`
	Amount    string `binding:"required" json:"amount"`
	Asset     string `                   json:"asset"`
	ID        int    `                   json:"id"`
}

type ReIssueRequest struct {
	Asset  string `binding:"required" json:"asset"`
	Amount string `binding:"required" json:"amount"`
	ID     int    `                   json:"id"`
}

type IssueMachineNFTRequest struct {
	Name           string `binding:"required" json:"name"`
	MachineAddress string `binding:"required" json:"machine-address"`
	Domain         string `binding:"required" json:"domain"`
	ID             int    `                   json:"id"`
}
