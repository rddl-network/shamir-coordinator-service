package backend

import "encoding/binary"

const (
	countKey      = "Count"
	taskKeyPrefix = "Task/"
)

func keyPrefix(p string) []byte {
	return []byte(p)
}

func taskKey(id int) []byte {
	var key []byte

	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(id))

	prefixBytes := []byte(taskKeyPrefix)
	key = append(key, prefixBytes...)
	key = append(key, buf...)

	return key
}
