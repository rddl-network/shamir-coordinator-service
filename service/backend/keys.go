package backend

import "encoding/binary"

const (
	countKeyPrefix          = "Count/"
	requestKeyPrefix        = "Request/"
	ReissueRequestPrefix    = "Reissue/"
	SendTokensRequestPrefix = "SendTokens/"
	IssueMachineNFTPrefix   = "IssueMachineNFTPrefix/"
)

func countKey(p string) []byte {
	return []byte(countKeyPrefix + p)
}

func requestKey(requestType string, id int) []byte {
	var key []byte

	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(id))
	key = append(key, []byte(requestType)...)
	key = append(key, buf...)

	return key
}
