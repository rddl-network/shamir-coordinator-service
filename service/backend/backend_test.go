package backend_test

import (
	"testing"

	"github.com/rddl-network/shamir-coordinator-service/service/backend"
	"github.com/rddl-network/shamir-coordinator-service/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/syndtr/goleveldb/leveldb"
)

func createNTasks(db *backend.DBConnector, n int) []backend.Task {
	items := make([]backend.Task, n)
	for i := range items {
		id, _ := db.CreateTask(items[i])
		items[i].ID = id
	}
	return items
}

func TestGetTask(t *testing.T) {
	db := testutil.SetupTestDBConnector(t)

	items := createNTasks(db, 1000)
	for _, item := range items {
		task, err := db.GetTask(item.ID)
		assert.NoError(t, err)
		assert.Equal(t, item, task)
	}
}

func TestGetAllTask(t *testing.T) {
	db := testutil.SetupTestDBConnector(t)

	items := createNTasks(db, 1000)
	tasks, err := db.GetAllTasks()
	assert.NoError(t, err)
	assert.Equal(t, items, tasks)
}

func TestDeleteTask(t *testing.T) {
	db := testutil.SetupTestDBConnector(t)

	items := createNTasks(db, 1000)

	db.DeleteTask(items[45].ID)
	db.DeleteTask(items[534].ID)

	tasks, err := db.GetAllTasks()
	assert.NoError(t, err)
	assert.Equal(t, len(tasks), 998)

	_, err = db.GetTask(items[45].ID)
	assert.Equal(t, leveldb.ErrNotFound, err)
}
