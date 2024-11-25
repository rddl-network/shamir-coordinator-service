package service_test

import (
	"sync"
	"testing"

	elements "github.com/rddl-network/elements-rpc"
	elementsmocks "github.com/rddl-network/elements-rpc/utils/mocks"
	"github.com/rddl-network/shamir-coordinator-service/testutil"
	"github.com/stretchr/testify/assert"
)

func TestLockingMutex(t *testing.T) {
	elements.Client = &elementsmocks.MockClient{}
	s := testutil.SetupTestService(t)

	passphrase := "password"
	err := s.PrepareWallet(passphrase)
	assert.NoError(t, err)

	err = s.PrepareWallet(passphrase)
	assert.NoError(t, err)

	locked, err := s.WalletLock()
	assert.NoError(t, err)
	assert.True(t, locked == false)

	locked, err = s.WalletLock()
	assert.NoError(t, err)
	assert.True(t, locked == true)

	err = s.PrepareWallet(passphrase)
	assert.NoError(t, err)

	locked, err = s.WalletLock()
	assert.NoError(t, err)
	assert.True(t, locked == true)
}

func TestSendTo(t *testing.T) {
	elements.Client = &elementsmocks.MockClient{}
	s := testutil.SetupTestService(t)

	address := "tlq1qqvsmfp0w3dmvwtkfteanzk0n7wksu6zx4pywzvak9p6d34yngghw39ynqwcxqrq3muqxffflmprr9exn8ldm79mlkz7dmpy0e"
	amount := "0.0001"
	txID, err := s.SendAsset(address, amount, "")
	assert.NoError(t, err)
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000000", txID)
}

func TestPrepareWallet(t *testing.T) {
	elements.Client = &elementsmocks.MockClient{}
	s := testutil.SetupTestService(t)

	passphrase := "password"
	err := s.PrepareWallet(passphrase)
	assert.NoError(t, err)
}

func TestReissueAsset(t *testing.T) {
	elements.Client = &elementsmocks.MockClient{}
	s := testutil.SetupTestService(t)

	asset := "06c20c8de513527f1ae6c901f74a05126525ac2d7e89306f4a7fd5ec4e674403"
	amount := "900.000"
	txID, err := s.ReissueAsset(asset, amount)
	assert.NoError(t, err)
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000000", txID)
}

func TestIssueNFTAsset(t *testing.T) {
	elements.Client = &elementsmocks.MockClient{}
	s := testutil.SetupTestService(t)

	var wg sync.WaitGroup

	for i := 0; i < 1; i++ {
		wg.Add(1)
		go func() {
			_, _, _, err := s.IssueNFTAsset("machine.Name", "machine.Address", "testnet-assets.rddl.io")
			assert.NoError(t, err)

			wg.Done()
		}()
	}
	wg.Wait()
}
