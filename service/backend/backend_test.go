package backend_test

import (
	"strconv"
	"testing"

	"github.com/rddl-network/shamir-coordinator-service/service/backend"
	"github.com/rddl-network/shamir-coordinator-service/testutil"
	"github.com/rddl-network/shamir-coordinator-service/types"
	"github.com/stretchr/testify/assert"
	"github.com/syndtr/goleveldb/leveldb"
)

func createNRequests(db *backend.DBConnector, requestType string, n int) []interface{} {
	items := make([]interface{}, n)
	for i := range items {
		iStr := strconv.Itoa(i)
		switch requestType {
		case backend.SendTokensRequestPrefix:
			items[i] = types.SendTokensRequest{
				Recipient: "recipient" + iStr,
				Amount:    "1000token",
				Asset:     "asset" + iStr,
			}
		case backend.ReissueRequestPrefix:
			items[i] = types.ReIssueRequest{
				Asset:  "asset" + iStr,
				Amount: "1000token",
			}
		case backend.IssueMachineNFTPrefix:
			items[i] = types.IssueMachineNFTRequest{
				Name:           "machine" + iStr,
				MachineAddress: "machAddr" + iStr,
				Domain:         "domain" + iStr,
			}
		}
		db.CreateRequest(requestType, items[i])
	}
	return items
}

func TestGetTask(t *testing.T) {
	db := testutil.SetupTestDBConnector(t)

	issueNFTitems := createNRequests(db, backend.IssueMachineNFTPrefix, 500)
	for i, item := range issueNFTitems {
		var request types.IssueMachineNFTRequest
		err := db.GetRequest(backend.IssueMachineNFTPrefix, i+1, &request)
		assert.NoError(t, err)
		assert.Equal(t, item, request)
	}

	reIssueItems := createNRequests(db, backend.ReissueRequestPrefix, 500)
	for i, item := range reIssueItems {
		var request types.ReIssueRequest
		err := db.GetRequest(backend.ReissueRequestPrefix, i+1, &request)
		assert.NoError(t, err)
		assert.Equal(t, item, request)
	}

	sendTokensItems := createNRequests(db, backend.SendTokensRequestPrefix, 500)
	for i, item := range sendTokensItems {
		var request types.SendTokensRequest
		err := db.GetRequest(backend.SendTokensRequestPrefix, i+1, &request)
		assert.NoError(t, err)
		assert.Equal(t, item, request)
	}
}

func TestGetAllTask(t *testing.T) {
	db := testutil.SetupTestDBConnector(t)

	sendTokensItems := createNRequests(db, backend.SendTokensRequestPrefix, 500)
	var comparableSendTokensItems []types.SendTokensRequest
	for _, r := range sendTokensItems {
		if req, ok := r.(types.SendTokensRequest); ok {
			comparableSendTokensItems = append(comparableSendTokensItems, req)
		}
	}
	sendTokensRequests, err := db.GetAllSendTokensRequests()
	assert.NoError(t, err)
	assert.Equal(t, comparableSendTokensItems, sendTokensRequests)

	reIssueItems := createNRequests(db, backend.ReissueRequestPrefix, 500)
	var comparableReIssueItems []types.ReIssueRequest
	for _, r := range reIssueItems {
		if req, ok := r.(types.ReIssueRequest); ok {
			comparableReIssueItems = append(comparableReIssueItems, req)
		}
	}
	reIssueRequests, err := db.GetAllReissueRequests()
	assert.NoError(t, err)
	assert.Equal(t, comparableReIssueItems, reIssueRequests)

	issueNFTitems := createNRequests(db, backend.IssueMachineNFTPrefix, 500)
	var comparableIssueNFTItems []types.IssueMachineNFTRequest
	for _, r := range issueNFTitems {
		if req, ok := r.(types.IssueMachineNFTRequest); ok {
			comparableIssueNFTItems = append(comparableIssueNFTItems, req)
		}
	}
	issueNFTRequests, err := db.GetAllIssueMachineNFTRequests()
	assert.NoError(t, err)
	assert.Equal(t, comparableIssueNFTItems, issueNFTRequests)
}

func TestDeleteTask(t *testing.T) {
	db := testutil.SetupTestDBConnector(t)

	createNRequests(db, backend.IssueMachineNFTPrefix, 500)
	createNRequests(db, backend.ReissueRequestPrefix, 500)
	createNRequests(db, backend.SendTokensRequestPrefix, 500)

	db.DeleteRequest(backend.IssueMachineNFTPrefix, 47+1)
	db.DeleteRequest(backend.ReissueRequestPrefix, 68+1)
	db.DeleteRequest(backend.ReissueRequestPrefix, 155+1)

	reqs, err := db.GetAllIssueMachineNFTRequests()
	assert.NoError(t, err)
	assert.Equal(t, len(reqs), 499)

	var request types.IssueMachineNFTRequest
	err = db.GetRequest(backend.IssueMachineNFTPrefix, 47+1, request)
	assert.Equal(t, leveldb.ErrNotFound, err)
}
