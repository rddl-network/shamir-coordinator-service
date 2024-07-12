package backend

import (
	"encoding/json"
	"errors"
	"strconv"
	"sync"

	"github.com/rddl-network/shamir-coordinator-service/config"
	"github.com/rddl-network/shamir-coordinator-service/types"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type DBConnector struct {
	db *leveldb.DB
}

var (
	dbMutex sync.Mutex
)

func NewDBConnector(db *leveldb.DB) *DBConnector {
	return &DBConnector{db: db}
}

func InitDB(cfg *config.Config) (db *leveldb.DB, err error) {
	return leveldb.OpenFile(cfg.DBPath, nil)
}

func (dc *DBConnector) IncrementCount(requestType string) (count int, err error) {
	countBytes, err := dc.db.Get(countKey(requestType), nil)
	if err != nil && !errors.Is(err, leveldb.ErrNotFound) {
		return 0, err
	}

	if countBytes == nil {
		count = 1
	} else {
		count, err = strconv.Atoi(string(countBytes))
		if err != nil {
			return 0, err
		}
		count++
	}

	dbMutex.Lock()
	defer dbMutex.Unlock()
	err = dc.db.Put(countKey(requestType), []byte(strconv.Itoa(count)), nil)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (dc *DBConnector) CreateRequest(requestType string, id int, request interface{}) (err error) {
	key := requestKey(requestType, id)
	val, err := json.Marshal(request)
	if err != nil {
		return err
	}

	dbMutex.Lock()
	defer dbMutex.Unlock()

	if err := dc.db.Put(key, val, nil); err != nil {
		return err
	}

	return nil
}

func (dc *DBConnector) GetRequest(requestType string, id int, request interface{}) (err error) {
	key := requestKey(requestType, id)
	dbMutex.Lock()
	defer dbMutex.Unlock()
	valBytes, err := dc.db.Get(key, nil)
	if err != nil {
		return
	}
	err = json.Unmarshal(valBytes, &request)
	return
}

func (dc *DBConnector) DeleteRequest(requestType string, id int) (err error) {
	key := requestKey(requestType, id)
	dbMutex.Lock()
	defer dbMutex.Unlock()
	return dc.db.Delete(key, nil)
}

func (dc *DBConnector) GetAllIssueMachineNFTRequests() (requests []types.IssueMachineNFTRequest, err error) {
	iter := dc.db.NewIterator(util.BytesPrefix([]byte(IssueMachineNFTPrefix)), nil)
	defer iter.Release()
	for iter.Next() {
		var request types.IssueMachineNFTRequest
		requestBytes := iter.Value()
		err = json.Unmarshal(requestBytes, &request)
		if err != nil {
			return
		}
		requests = append(requests, request)
	}
	return
}

func (dc *DBConnector) GetAllSendTokensRequests() (requests []types.SendTokensRequest, err error) {
	iter := dc.db.NewIterator(util.BytesPrefix([]byte(SendTokensRequestPrefix)), nil)
	defer iter.Release()
	for iter.Next() {
		var request types.SendTokensRequest
		requestBytes := iter.Value()
		err = json.Unmarshal(requestBytes, &request)
		if err != nil {
			return
		}
		requests = append(requests, request)
	}
	return
}

func (dc *DBConnector) GetAllReissueRequests() (requests []types.ReIssueRequest, err error) {
	iter := dc.db.NewIterator(util.BytesPrefix([]byte(SendTokensRequestPrefix)), nil)
	defer iter.Release()
	for iter.Next() {
		var request types.ReIssueRequest
		requestBytes := iter.Value()
		err = json.Unmarshal(requestBytes, &request)
		if err != nil {
			return
		}
		requests = append(requests, request)
	}
	return
}
