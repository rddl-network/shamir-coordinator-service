package service_test

import (
	"testing"

	"github.com/rddl-network/shamir-coordinator-service/testutil"
	"github.com/stretchr/testify/assert"
)

func TestCollectMnemonics(t *testing.T) {
	s := testutil.SetupTestService(t)

	mnemonics, err := s.CollectMnemonics()
	assert.NoError(t, err)
	assert.Equal(t, 3, len(mnemonics))
}
