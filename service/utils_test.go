package service_test

import (
	"testing"

	hexutil "github.com/rddl-network/go-utils/hex"
	strutil "github.com/rddl-network/go-utils/str"
	"github.com/stretchr/testify/assert"
)

func TestIsValidHex(t *testing.T) {
	assert.False(t, hexutil.IsValidHex("31622fc2d536a751dfff93c6cf21b3d206d4c5362f7fa48e974233db0a56c6c73e6c7466e424d1fd04ed5e0e94e155a"))
	assert.True(t, hexutil.IsValidHex("31622fc2d536a751dfff93c6cf21b3d206d4c5362f7fa48e974233db0a56c6c73e6c7466e424d1fd04ed5e0e94e155ad"))
}

func TestContainsString(t *testing.T) {
	array := []string{"banana", "apple", "tree"}
	assert.True(t, strutil.ContainsString(array, "tree"))
	assert.False(t, strutil.ContainsString(array, "flower"))
}
