// Package testutils provides shared test helpers for taskadmin:
// an in-memory SQLite task store and a mock Layout implementation.
package testutils

import (
	"database/sql"
	"net/http"
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

// MockLayout is a test implementation of the shared.Layout interface.
// It records the values passed to its setters for assertions.
type MockLayout struct {
	Title       string
	ScriptURLs  []string
	Scripts     []string
	StyleURLs   []string
	Styles      []string
	Body        string
	RenderCount int
}

// NewMockLayout creates a new mock layout for tests.
func NewMockLayout() *MockLayout {
	return &MockLayout{}
}

// SetTitle records the title
func (m *MockLayout) SetTitle(title string) { m.Title = title }

// SetScriptURLs records the script URLs
func (m *MockLayout) SetScriptURLs(scripts []string) { m.ScriptURLs = scripts }

// SetScripts records the scripts
func (m *MockLayout) SetScripts(scripts []string) { m.Scripts = scripts }

// SetStyleURLs records the style URLs
func (m *MockLayout) SetStyleURLs(styles []string) { m.StyleURLs = styles }

// SetStyles records the styles
func (m *MockLayout) SetStyles(styles []string) { m.Styles = styles }

// SetBody records the body
func (m *MockLayout) SetBody(body string) { m.Body = body }

// Render returns the recorded body and increments the render count
func (m *MockLayout) Render(w http.ResponseWriter, r *http.Request) string {
	m.RenderCount++
	return m.Body
}
