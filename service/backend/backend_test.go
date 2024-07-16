package backend_test

import (
	"strconv"
	"testing"

	"github.com/rddl-network/shamir-coordinator-service/service/backend"
	"github.com/rddl-network/shamir-coordinator-service/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/syndtr/goleveldb/leveldb"
)

func createNRequests(db *backend.DBConnector, requestType string, n int) []interface{} {
	items := make([]interface{}, n)
	for i := range items {
		iStr := strconv.Itoa(i)
		id, _ := db.IncrementCount(requestType)
		switch requestType {
		case backend.SendTokensRequestPrefix:
			items[i] = backend.SendTokensRequest{
				Recipient: "recipient" + iStr,
				Amount:    "1000",
				Asset:     "asset" + iStr,
				ID:        id,
			}
		case backend.ReissueRequestPrefix:
			items[i] = backend.ReIssueRequest{
				Asset:  "asset" + iStr,
				Amount: "1500",
				ID:     id,
			}
		case backend.IssueMachineNFTPrefix:
			items[i] = backend.IssueMachineNFTRequest{
				Name:           "machine" + iStr,
				MachineAddress: "machAddr" + iStr,
				Domain:         "domain" + iStr,
				ID:             id,
			}
		}
		_ = db.CreateRequest(requestType, id, items[i])
	}
	return items
}

func TestGetRequests(t *testing.T) {
	db := testutil.SetupTestDBConnector(t)

	issueNFTitems := createNRequests(db, backend.IssueMachineNFTPrefix, 400)
	for i, item := range issueNFTitems {
		var request backend.IssueMachineNFTRequest
		err := db.GetRequest(backend.IssueMachineNFTPrefix, i+1, &request)
		assert.NoError(t, err)
		assert.Equal(t, item, request)
	}

	reIssueItems := createNRequests(db, backend.ReissueRequestPrefix, 500)
	for i, item := range reIssueItems {
		var request backend.ReIssueRequest
		err := db.GetRequest(backend.ReissueRequestPrefix, i+1, &request)
		assert.NoError(t, err)
		assert.Equal(t, item, request)
	}

	sendTokensItems := createNRequests(db, backend.SendTokensRequestPrefix, 300)
	for i, item := range sendTokensItems {
		var request backend.SendTokensRequest
		err := db.GetRequest(backend.SendTokensRequestPrefix, i+1, &request)
		assert.NoError(t, err)
		assert.Equal(t, item, request)
	}
}

func TestGetAllRequests(t *testing.T) {
	db := testutil.SetupTestDBConnector(t)

	sendTokensItems := createNRequests(db, backend.SendTokensRequestPrefix, 500)
	var comparableSendTokensItems []backend.SendTokensRequest
	for _, r := range sendTokensItems {
		if req, ok := r.(backend.SendTokensRequest); ok {
			comparableSendTokensItems = append(comparableSendTokensItems, req)
		}
	}
	sendTokensRequests, err := db.GetAllSendTokensRequests()
	assert.NoError(t, err)
	assert.Equal(t, comparableSendTokensItems, sendTokensRequests)

	reIssueItems := createNRequests(db, backend.ReissueRequestPrefix, 500)
	var comparableReIssueItems []backend.ReIssueRequest
	for _, r := range reIssueItems {
		if req, ok := r.(backend.ReIssueRequest); ok {
			comparableReIssueItems = append(comparableReIssueItems, req)
		}
	}
	reIssueRequests, err := db.GetAllReissueRequests()
	assert.NoError(t, err)
	assert.Equal(t, comparableReIssueItems, reIssueRequests)

	issueNFTitems := createNRequests(db, backend.IssueMachineNFTPrefix, 500)
	var comparableIssueNFTItems []backend.IssueMachineNFTRequest
	for _, r := range issueNFTitems {
		if req, ok := r.(backend.IssueMachineNFTRequest); ok {
			comparableIssueNFTItems = append(comparableIssueNFTItems, req)
		}
	}
	issueNFTRequests, err := db.GetAllIssueMachineNFTRequests()
	assert.NoError(t, err)
	assert.Equal(t, comparableIssueNFTItems, issueNFTRequests)
}

func TestDeleteRequest(t *testing.T) {
	db := testutil.SetupTestDBConnector(t)

	nftReqs := createNRequests(db, backend.IssueMachineNFTPrefix, 500)
	for _, r := range nftReqs {
		if req, ok := r.(backend.IssueMachineNFTRequest); ok {
			err := db.DeleteRequest(backend.IssueMachineNFTPrefix, req.ID)
			assert.NoError(t, err)

			var request backend.IssueMachineNFTRequest
			err = db.GetRequest(backend.IssueMachineNFTPrefix, req.ID, request)
			assert.Equal(t, leveldb.ErrNotFound, err)
		}
	}
	nReqs, err := db.GetAllIssueMachineNFTRequests()
	assert.NoError(t, err)
	assert.Equal(t, len(nReqs), 0)

	reIssueReqs := createNRequests(db, backend.ReissueRequestPrefix, 500)
	for _, r := range reIssueReqs {
		if req, ok := r.(backend.ReIssueRequest); ok {
			err := db.DeleteRequest(backend.ReissueRequestPrefix, req.ID)
			assert.NoError(t, err)

			var request backend.ReIssueRequest
			err = db.GetRequest(backend.ReissueRequestPrefix, req.ID, request)
			assert.Equal(t, leveldb.ErrNotFound, err)
		}
	}
	rReqs, err := db.GetAllReissueRequests()
	assert.NoError(t, err)
	assert.Equal(t, len(rReqs), 0)

	sendTokensReqs := createNRequests(db, backend.SendTokensRequestPrefix, 500)
	for _, r := range sendTokensReqs {
		if req, ok := r.(backend.SendTokensRequest); ok {
			err := db.DeleteRequest(backend.SendTokensRequestPrefix, req.ID)
			assert.NoError(t, err)

			var request backend.SendTokensRequest
			err = db.GetRequest(backend.SendTokensRequestPrefix, req.ID, request)
			assert.Equal(t, leveldb.ErrNotFound, err)
		}
	}
	sReqs, err := db.GetAllSendTokensRequests()
	assert.NoError(t, err)
	assert.Equal(t, len(sReqs), 0)
}
