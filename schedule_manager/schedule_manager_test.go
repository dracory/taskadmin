package schedule_manager_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/dracory/taskadmin/schedule_manager"
	"github.com/dracory/taskadmin/shared"
	"github.com/dracory/taskadmin/testutils"
	"github.com/dracory/taskstore"
)

type testEnv struct {
	ui    schedule_manager.UiInterface
	store taskstore.StoreInterface
}

func newTestEnv(t *testing.T) *testEnv {
	t.Helper()

	store := testutils.NewStore(t)

	ui := schedule_manager.UI(shared.UiConfig{
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

	r := httptest.NewRequest("POST", "/admin/task?controller=schedule-manager&action="+action, bytes.NewReader(payload))
	w := httptest.NewRecorder()

	e.ui.ScheduleManager(w, r)
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

func Test_Page_renders(t *testing.T) {
	env := newTestEnv(t)

	r := httptest.NewRequest("GET", "/admin/task?controller=schedule-manager", nil)
	w := httptest.NewRecorder()

	env.ui.ScheduleManager(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "Schedules") {
		t.Error("expected page to render schedules content")
	}
}

func Test_Create_and_list_schedule(t *testing.T) {
	env := newTestEnv(t)
	definition := env.seedDefinition(t, "test-task")

	w := env.post(t, "create-schedule", map[string]any{
		"name":               "Nightly Report",
		"status":             "active",
		"queue_name":         "default",
		"task_definition_id": definition.GetID(),
		"start_at":           "2099-01-01T00:00",
		"frequency":          "daily",
		"interval":           1,
	})
	if w.Code != http.StatusOK {
		t.Fatalf("create: expected 200, got %d", w.Code)
	}

	w = env.post(t, "load-schedules", map[string]any{})
	var result struct {
		Status string `json:"status"`
		Data   struct {
			Schedules []map[string]any `json:"schedules"`
			Total     int              `json:"total"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("list: invalid json: %v", err)
	}
	if result.Data.Total != 1 || len(result.Data.Schedules) != 1 {
		t.Fatalf("expected 1 schedule, got %d", result.Data.Total)
	}

	created := result.Data.Schedules[0]
	if created["name"] != "Nightly Report" {
		t.Errorf("expected name 'Nightly Report', got %v", created["name"])
	}
	if created["next_run_at"] == "" || created["next_run_at"] == nil {
		t.Error("expected next_run_at to be calculated")
	}
}

func Test_Create_requires_name_and_definition(t *testing.T) {
	env := newTestEnv(t)

	w := env.post(t, "create-schedule", map[string]any{"name": ""})
	if !strings.Contains(w.Body.String(), "Name is required") {
		t.Errorf("expected name required error, got %s", w.Body.String())
	}

	w = env.post(t, "create-schedule", map[string]any{"name": "x"})
	if !strings.Contains(w.Body.String(), "Task definition is required") {
		t.Errorf("expected task definition required error, got %s", w.Body.String())
	}
}

func Test_Update_and_delete_schedule(t *testing.T) {
	env := newTestEnv(t)
	definition := env.seedDefinition(t, "test-task")

	schedule := taskstore.NewSchedule().
		SetName("Old Name").
		SetStatus("draft").
		SetQueueName("default").
		SetTaskDefinitionID(definition.GetID())
	if err := env.store.ScheduleCreate(t.Context(), schedule); err != nil {
		t.Fatalf("seed schedule: %v", err)
	}

	w := env.post(t, "update-schedule", map[string]any{
		"id":                 schedule.GetID(),
		"name":               "New Name",
		"status":             "active",
		"queue_name":         "default",
		"task_definition_id": definition.GetID(),
		"frequency":          "hourly",
		"interval":           2,
	})
	if !strings.Contains(w.Body.String(), "success") {
		t.Fatalf("update failed: %s", w.Body.String())
	}

	updated, err := env.store.ScheduleFindByID(t.Context(), schedule.GetID())
	if err != nil || updated == nil {
		t.Fatalf("find updated: %v", err)
	}
	if updated.GetName() != "New Name" {
		t.Errorf("expected 'New Name', got %s", updated.GetName())
	}
	if updated.GetRecurrenceRule().GetFrequency() != taskstore.FrequencyHourly {
		t.Errorf("expected hourly frequency, got %s", updated.GetRecurrenceRule().GetFrequency())
	}

	w = env.post(t, "delete-schedule", map[string]any{"id": schedule.GetID()})
	if !strings.Contains(w.Body.String(), "success") {
		t.Fatalf("delete failed: %s", w.Body.String())
	}

	list, err := env.store.ScheduleList(t.Context(), taskstore.NewScheduleQuery().SetLimit(0))
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(list) != 0 {
		t.Errorf("expected schedule to be soft deleted, %d remain", len(list))
	}
}
