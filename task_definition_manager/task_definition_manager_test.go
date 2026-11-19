package task_definition_manager_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/dracory/taskadmin/shared"
	"github.com/dracory/taskadmin/task_definition_manager"
	"github.com/dracory/taskadmin/testutils"
	"github.com/dracory/taskstore"
)

type testEnv struct {
	ui    task_definition_manager.UiInterface
	store taskstore.StoreInterface
}

func newTestEnv(t *testing.T) *testEnv {
	t.Helper()

	store := testutils.NewStore(t)

	ui := task_definition_manager.UI(shared.UiConfig{
		Store:  store,
		Logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
		Layout: func(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options shared.LayoutOptions) string {
			return "<html><head><title>" + webpageTitle + "</title></head><body>" + webpageHtml + "</body></html>"
		},
	})

	return &testEnv{ui: ui, store: store}
}

func (e *testEnv) post(t *testing.T, action string, body map[string]any) *httptest.ResponseRecorder {
	t.Helper()

	payload, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest("POST", "/admin/task?controller=task-definition-manager&action="+action, bytes.NewReader(payload))
	w := httptest.NewRecorder()

	e.ui.TaskDefinitionManager(w, r)
	return w
}

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

func Test_UI(t *testing.T) {
	env := newTestEnv(t)
	if env.ui == nil {
		t.Fatal("UI() should return a non-nil controller")
	}
	if env.ui.Store() == nil {
		t.Fatal("UI() should set the store")
	}
}

func Test_renderPage(t *testing.T) {
	env := newTestEnv(t)

	r := httptest.NewRequest("GET", "/admin/task?controller=task-definition-manager", nil)
	w := httptest.NewRecorder()
	env.ui.TaskDefinitionManager(w, r)

	if w.Code != 200 {
		t.Fatalf("page should return 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "Tasks. Definitions") {
		t.Error("page should render the definitions header")
	}
}

func Test_loadDefinitions(t *testing.T) {
	env := newTestEnv(t)
	env.seedDefinition(t, "load-me")

	w := env.post(t, "load-definitions", map[string]any{"page": 0, "per_page": 10})

	var resp struct {
		Status string `json:"status"`
		Data   struct {
			Definitions []map[string]any `json:"definitions"`
			Total       int              `json:"total"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}
	if resp.Status != "success" || resp.Data.Total != 1 {
		t.Fatalf("expected 1 definition, got %s", w.Body.String())
	}
}

func Test_definitionCreate(t *testing.T) {
	env := newTestEnv(t)

	w := env.post(t, "create-definition", map[string]any{
		"title":       "My Task",
		"alias":       "my-task",
		"status":      "active",
		"description": "desc",
	})
	if !strings.Contains(w.Body.String(), "success") {
		t.Fatalf("expected success, got %s", w.Body.String())
	}

	found, err := env.store.TaskDefinitionFindByAlias(t.Context(), "my-task")
	if err != nil || found == nil {
		t.Fatal("expected definition to be created")
	}
}

func Test_definitionUpdate(t *testing.T) {
	env := newTestEnv(t)
	definition := env.seedDefinition(t, "update-me")

	w := env.post(t, "update-definition", map[string]any{
		"id":          definition.GetID(),
		"title":       "Renamed",
		"alias":       "update-me",
		"status":      "inactive",
		"description": "new desc",
	})
	if !strings.Contains(w.Body.String(), "success") {
		t.Fatalf("expected success, got %s", w.Body.String())
	}

	reloaded, err := env.store.TaskDefinitionFindByID(t.Context(), definition.GetID())
	if err != nil {
		t.Fatal(err)
	}
	if reloaded.GetTitle() != "Renamed" || reloaded.GetStatus() != "inactive" {
		t.Errorf("expected updated fields, got %s / %s", reloaded.GetTitle(), reloaded.GetStatus())
	}
}

func Test_definitionDelete(t *testing.T) {
	env := newTestEnv(t)
	definition := env.seedDefinition(t, "delete-me")

	w := env.post(t, "delete-definition", map[string]any{"id": definition.GetID()})
	if !strings.Contains(w.Body.String(), "success") {
		t.Fatalf("expected success, got %s", w.Body.String())
	}

	count, err := env.store.TaskDefinitionCount(t.Context(), taskstore.TaskDefinitionQuery())
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Errorf("expected 0 definitions after delete, got %d", count)
	}
}
