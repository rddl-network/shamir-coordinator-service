package types

type SendTokensRequest struct {
	Recipient string `binding:"required" json:"recipient"`
	Amount    string `binding:"required" json:"amount"`
	Asset     string `                   json:"asset"`
}

type SendTokensResponse struct {
	TxID string `binding:"required" json:"tx-id"`
}

type ReIssueRequest struct {
	Asset  string `binding:"required" json:"asset"`
	Amount string `binding:"required" json:"amount"`
}

type ReIssueResponse struct {
	TxID string `binding:"required" json:"tx-id"`
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
