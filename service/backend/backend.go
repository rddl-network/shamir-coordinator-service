package backend

import (
	"encoding/json"
	"errors"
	"strconv"
	"sync"

	"github.com/rddl-network/shamir-coordinator-service/config"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type Task struct {
	ID int `binding:"gte=0" json:"id"`
}

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

func (dc *DBConnector) incrementCount() (count int, err error) {
	countBytes, err := dc.db.Get(keyPrefix(countKey), nil)
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
	err = dc.db.Put(keyPrefix(countKey), []byte(strconv.Itoa(count)), nil)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (dc *DBConnector) CreateTask(task Task) (id int, err error) {
	id, err = dc.incrementCount()
	if err != nil {
		return
	}

	task.ID = id

	key := taskKey(id)
	val, err := json.Marshal(task)
	if err != nil {
		return 0, err
	}

	dbMutex.Lock()
	defer dbMutex.Unlock()

	if err := dc.db.Put(key, val, nil); err != nil {
		return 0, err
	}

	return id, nil
}

func (dc *DBConnector) GetTask(id int) (task Task, err error) {
	key := taskKey(id)
	dbMutex.Lock()
	defer dbMutex.Unlock()
	valBytes, err := dc.db.Get(key, nil)
	if err != nil {
		return
	}
	err = json.Unmarshal(valBytes, &task)
	return
}

func (dc *DBConnector) DeleteTask(id int) (err error) {
	key := taskKey(id)
	dbMutex.Lock()
	defer dbMutex.Unlock()
	return dc.db.Delete(key, nil)
}

func (dc *DBConnector) GetAllTasks() (tasks []Task, err error) {
	iter := dc.db.NewIterator(util.BytesPrefix([]byte(taskKeyPrefix)), nil)
	defer iter.Release()
	for iter.Next() {
		var task Task
		taskBytes := iter.Value()
		err = json.Unmarshal(taskBytes, &task)
		if err != nil {
			return
		}
		tasks = append(tasks, task)
	}
	return
}
