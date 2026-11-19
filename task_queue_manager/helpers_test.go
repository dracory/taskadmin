package task_queue_manager_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/dracory/taskadmin/shared"
	"github.com/dracory/taskadmin/task_queue_manager"
	"github.com/dracory/taskadmin/testutils"
	"github.com/dracory/taskstore"
)

// testEnv bundles a task queue manager controller UI backed by an
// in-memory SQLite taskstore, ready to serve HTTP requests in a test.
type testEnv struct {
	ui    task_queue_manager.UiInterface
	store taskstore.StoreInterface
}

// newTestEnv boots an in-memory SQLite taskstore and a task queue
// manager controller UI.
func newTestEnv(t *testing.T) *testEnv {
	t.Helper()

	store := testutils.NewStore(t)

	ui := task_queue_manager.UI(shared.UiConfig{
		Store:  store,
		Logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
		Layout: func(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options shared.LayoutOptions) string {
			return "<html><head><title>" + webpageTitle + "</title></head><body>" + webpageHtml + "</body></html>"
		},
	})

	return &testEnv{ui: ui, store: store}
}

// seedDefinition creates an active task definition with the given alias.
func (e *testEnv) seedDefinition(t *testing.T, alias string) taskstore.TaskDefinitionInterface {
	t.Helper()

	definition := taskstore.NewTaskDefinition().
		SetTitle("Definition " + alias).
		SetAlias(alias).
		SetStatus(taskstore.TaskDefinitionStatusActive)

	if err := e.store.TaskDefinitionCreate(t.Context(), definition); err != nil {
		t.Fatalf("seed definition: %v", err)
	}

	found, err := e.store.TaskDefinitionFindByAlias(t.Context(), alias)
	if err != nil || found == nil {
		t.Fatalf("seed definition lookup: %v", err)
	}
	return found
}

// seedQueueTask enqueues a task for the given definition alias.
func (e *testEnv) seedQueueTask(t *testing.T, alias string) taskstore.TaskQueueInterface {
	t.Helper()

	queued, err := e.store.TaskDefinitionEnqueueByAlias(t.Context(), taskstore.DefaultQueueName, alias, map[string]any{"k": "v"})
	if err != nil {
		t.Fatalf("seed queue task: %v", err)
	}
	return queued
}

// post sends a POST JSON request to the controller's TaskQueueManager
// entrypoint with the given action and body, returning the recorder.
func (e *testEnv) post(t *testing.T, action string, body map[string]any) *httptest.ResponseRecorder {
	t.Helper()

	payload, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest("POST", "/admin/task?controller=task-queue-manager&action="+action, bytes.NewReader(payload))
	w := httptest.NewRecorder()

	e.ui.TaskQueueManager(w, r)
	return w
}

// get renders the page with no action
func (e *testEnv) get(t *testing.T) *httptest.ResponseRecorder {
	t.Helper()

	r := httptest.NewRequest("GET", "/admin/task?controller=task-queue-manager", nil)
	w := httptest.NewRecorder()

	e.ui.TaskQueueManager(w, r)
	return w
}
