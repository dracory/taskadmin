// Package testutils provides shared test helpers for taskadmin:
// an in-memory SQLite task store and a mock Layout implementation.
package testutils

import (
	"database/sql"
	"testing"

	"github.com/dracory/taskstore"
	_ "modernc.org/sqlite"
)

// NewStore creates a real taskstore backed by an in-memory SQLite database.
func NewStore(t *testing.T) taskstore.StoreInterface {
	t.Helper()

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	store, err := taskstore.NewStore(taskstore.NewStoreOptions{
		TaskDefinitionTableName: "task_definitions",
		TaskQueueTableName:      "task_queue",
		ScheduleTableName:       "schedules",
		DB:                      db,
		AutomigrateEnabled:      true,
	})
	if err != nil {
		t.Fatal(err)
	}

	return store
}
