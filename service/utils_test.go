package service_test

import (
	"testing"

	"github.com/rddl-network/shamir-coordinator-service/service"
	"github.com/stretchr/testify/assert"
)

func TestIsValidHex(t *testing.T) {
	assert.False(t, service.IsValidHex("31622fc2d536a751dfff93c6cf21b3d206d4c5362f7fa48e974233db0a56c6c73e6c7466e424d1fd04ed5e0e94e155a"))
	assert.True(t, service.IsValidHex("31622fc2d536a751dfff93c6cf21b3d206d4c5362f7fa48e974233db0a56c6c73e6c7466e424d1fd04ed5e0e94e155ad"))
}

func TestContainsString(t *testing.T) {
	array := []string{"banana", "apple", "tree"}
	assert.True(t, service.ContainsString(array, "tree"))
	assert.False(t, service.ContainsString(array, "flower"))
}
