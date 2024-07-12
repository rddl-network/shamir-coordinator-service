package types

type SendTokensRequest struct {
	Recipient string `binding:"required" json:"recipient"`
	Amount    string `binding:"required" json:"amount"`
	Asset     string `json:"asset"`
	ID        int    `json:"id"`
}

type SendTokensResponse struct {
	TxID string `binding:"required" json:"tx-id"`
}

type ReIssueRequest struct {
	Asset  string `binding:"required" json:"asset"`
	Amount string `binding:"required" json:"amount"`
	ID     int    `json:"id"`
}

type ReIssueResponse struct {
	TxID string `binding:"required" json:"tx-id"`
}

type IssueMachineNFTRequest struct {
	Name           string `binding:"required" json:"name"`
	MachineAddress string `binding:"required" json:"machine-address"`
	Domain         string `binding:"required" json:"domain"`
	ID             int    `json:"id"`
}

type IssueMachineNFTResponse struct {
	Asset    string `binding:"required" json:"asset"`
	Contract string `binding:"required" json:"contract"`
	HexTX    string `binding:"required" json:"hex-tx"`
}

type MnemonicsResponse struct {
	Mnemonics []string `binding:"required" json:"mnemonics"`
	Seed      string   `binding:"required" json:"seed"`
}

type Entity struct {
	Domain string `json:"domain"`
}

type Contract struct {
	Entity       Entity `json:"entity"`
	IssuerPubkey string `json:"issuer_pubkey"` //nolint:tagliatelle // the format liquid network needs it
	MachineAddr  string `json:"machine_addr"`  //nolint:tagliatelle // the format liquid network needs it
	Name         string `json:"name"`
	Precision    uint64 `json:"precision"`
	Version      uint64 `json:"version"`
}
