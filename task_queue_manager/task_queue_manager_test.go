package task_queue_manager_test

import (
	"encoding/json"
	"github.com/dracory/taskstore"
	"strings"
	"testing"
)

func Test_UI(t *testing.T) {
	env := newTestEnv(t)
	if env.ui == nil {
		t.Fatal("UI() should return a non-nil controller")
	}
	if env.ui.Store() == nil {
		t.Fatal("UI() should set the store")
	}
	if env.ui.Logger() == nil {
		t.Fatal("UI() should set the logger")
	}
}

func Test_renderPage(t *testing.T) {
	env := newTestEnv(t)
	w := env.get(t)

	if w.Code != 200 {
		t.Fatalf("page should return 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "Tasks. Queue") {
		t.Error("page should render the queue manager header")
	}
}

func Test_loadTasks(t *testing.T) {
	env := newTestEnv(t)
	env.seedDefinition(t, "test-task")
	env.seedQueueTask(t, "test-task")

	w := env.post(t, "load-tasks", map[string]any{
		"page":     0,
		"per_page": 10,
	})

	var resp struct {
		Status string `json:"status"`
		Data   struct {
			Tasks []map[string]any `json:"tasks"`
			Total int              `json:"total"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}
	if resp.Status != "success" {
		t.Fatalf("expected success, got %s", resp.Status)
	}
	if resp.Data.Total != 1 || len(resp.Data.Tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(resp.Data.Tasks))
	}
}

func Test_enqueueTask(t *testing.T) {
	env := newTestEnv(t)
	definition := env.seedDefinition(t, "enqueue-me")

	w := env.post(t, "enqueue-task", map[string]any{
		"task_definition_id": definition.GetID(),
		"parameters":         `{"a":"b"}`,
	})

	var resp struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}
	if resp.Status != "success" {
		t.Fatalf("expected success, got %s", resp.Status)
	}

	count, err := env.store.TaskQueueCount(t.Context(), taskstore.TaskQueueQuery())
	if err != nil || count != 1 {
		t.Fatalf("expected 1 queued task, got %d (err %v)", count, err)
	}
}

func Test_enqueueTask_invalidJSON(t *testing.T) {
	env := newTestEnv(t)
	definition := env.seedDefinition(t, "bad-json")

	w := env.post(t, "enqueue-task", map[string]any{
		"task_definition_id": definition.GetID(),
		"parameters":         "not json",
	})

	if !strings.Contains(w.Body.String(), "error") {
		t.Error("expected error for invalid JSON parameters")
	}
}

func Test_taskDetails(t *testing.T) {
	env := newTestEnv(t)
	env.seedDefinition(t, "detail-task")
	queued := env.seedQueueTask(t, "detail-task")

	w := env.post(t, "task-details", map[string]any{"id": queued.GetID()})

	var resp struct {
		Status string `json:"status"`
		Data   struct {
			ID         string `json:"id"`
			Status     string `json:"status"`
			Parameters string `json:"parameters"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}
	if resp.Status != "success" || resp.Data.ID != queued.GetID() {
		t.Fatalf("expected task details, got %s", w.Body.String())
	}
}

func Test_taskRestart(t *testing.T) {
	env := newTestEnv(t)
	env.seedDefinition(t, "restart-task")
	queued := env.seedQueueTask(t, "restart-task")

	queued.SetStatus("failed")
	if err := env.store.TaskQueueUpdate(t.Context(), queued); err != nil {
		t.Fatal(err)
	}

	w := env.post(t, "restart-task", map[string]any{"id": queued.GetID()})
	if !strings.Contains(w.Body.String(), "success") {
		t.Fatalf("expected success, got %s", w.Body.String())
	}

	reloaded, err := env.store.TaskQueueFindByID(t.Context(), queued.GetID())
	if err != nil {
		t.Fatal(err)
	}
	if reloaded.GetStatus() != "queued" {
		t.Errorf("expected status queued, got %s", reloaded.GetStatus())
	}
}

func Test_taskDelete(t *testing.T) {
	env := newTestEnv(t)
	env.seedDefinition(t, "delete-task")
	queued := env.seedQueueTask(t, "delete-task")

	w := env.post(t, "delete-task", map[string]any{"id": queued.GetID()})
	if !strings.Contains(w.Body.String(), "success") {
		t.Fatalf("expected success, got %s", w.Body.String())
	}

	count, err := env.store.TaskQueueCount(t.Context(), taskstore.TaskQueueQuery())
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Errorf("expected 0 tasks after delete, got %d", count)
	}
}

func Test_taskUpdateParams(t *testing.T) {
	env := newTestEnv(t)
	env.seedDefinition(t, "params-task")
	queued := env.seedQueueTask(t, "params-task")

	w := env.post(t, "update-parameters", map[string]any{
		"id":         queued.GetID(),
		"parameters": `{"x":"y"}`,
	})
	if !strings.Contains(w.Body.String(), "success") {
		t.Fatalf("expected success, got %s", w.Body.String())
	}

	reloaded, err := env.store.TaskQueueFindByID(t.Context(), queued.GetID())
	if err != nil {
		t.Fatal(err)
	}
	if reloaded.GetParameters() != `{"x":"y"}` {
		t.Errorf("expected updated parameters, got %s", reloaded.GetParameters())
	}
}
