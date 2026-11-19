package task_queue_task_restart

import (
	"github.com/dracory/taskadmin/shared"
	"github.com/dracory/taskadmin/testutils"
	"github.com/dracory/taskstore"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_taskQueueTaskRestart(t *testing.T) {
	// Test taskQueueTaskRestart controller constructor with real SQLite store
	store := testutils.NewStore(t)

	controller := testUI(store)

	if controller == nil {
		t.Error("taskQueueTaskRestart() should return a non-nil controller")
	}
	if controller.Store() == nil {
		t.Error("taskQueueTaskRestart() should set store")
	}
}

func Test_taskQueueTaskRestartController_struct_fields(t *testing.T) {
	// Test ui struct fields
	store := testutils.NewStore(t)

	controller := testUI(store)

	if controller.Logger() == nil {
		t.Error("ui logger should be set")
	}
	if controller.Store() == nil {
		t.Error("ui store should not be nil")
	}
}

func Test_taskQueueTaskRestartController_with_nil_store(t *testing.T) {
	// Test ui with nil store

	controller := testUI(nil)

	if controller.Store() != nil {
		t.Error("ui store should be nil")
	}
}

func Test_taskQueueTaskRestartControllerData_request_field(t *testing.T) {
	// Test taskQueueTaskRestartControllerData request field
	req := httptest.NewRequest("GET", "/", nil)
	data := taskQueueTaskRestartControllerData{request: req}

	if data.request == nil {
		t.Error("taskQueueTaskRestartControllerData request should not be nil")
	}
}

func Test_taskQueueTaskRestartControllerData_nil_request(t *testing.T) {
	// Test taskQueueTaskRestartControllerData with nil request
	data := taskQueueTaskRestartControllerData{}

	if data.request != nil {
		t.Error("taskQueueTaskRestartControllerData request should be nil")
	}
}

// testUI builds a controller instance for tests
func testUI(store taskstore.StoreInterface) *ui {
	return &ui{UiBase: shared.NewUiBase(shared.UiConfig{
		Store:  store,
		Logger: slog.Default(),
		FuncLayout: func(r *http.Request) shared.Layout {
			return testutils.NewMockLayout()
		},
	})}
}
