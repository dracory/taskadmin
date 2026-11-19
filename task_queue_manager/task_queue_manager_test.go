package task_queue_manager

import (
	"github.com/dracory/taskadmin/shared"
	"github.com/dracory/taskadmin/testutils"
	"github.com/dracory/taskstore"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_taskQueueManager(t *testing.T) {
	// Test taskQueueManager controller constructor with real SQLite store
	store := testutils.NewStore(t)

	controller := testUI(store)

	if controller == nil {
		t.Error("taskQueueManager() should return a non-nil controller")
	}
	if controller.Store() == nil {
		t.Error("taskQueueManager() should set store")
	}
	if controller.Layout(nil) == nil {
		t.Error("taskQueueManager() should set layout")
	}
}

func Test_taskQueueManagerController_struct_fields(t *testing.T) {
	// Test ui struct fields
	store := testutils.NewStore(t)

	controller := testUI(store)

	if controller.Logger() == nil {
		t.Error("ui logger should be set")
	}
	if controller.Store() == nil {
		t.Error("ui store should not be nil")
	}
	if controller.Layout(nil) == nil {
		t.Error("ui layout should not be nil")
	}
}

func Test_taskQueueManagerController_with_nil_store(t *testing.T) {
	// Test ui with nil store

	controller := testUI(nil)

	if controller.Store() != nil {
		t.Error("ui store should be nil")
	}
}

func Test_taskQueueManagerControllerData_request_field(t *testing.T) {
	// Test taskQueueManagerControllerData request field
	req := httptest.NewRequest("GET", "/", nil)
	data := taskQueueManagerControllerData{request: req}

	if data.request == nil {
		t.Error("taskQueueManagerControllerData request should not be nil")
	}
}

func Test_taskQueueManagerControllerData_nil_request(t *testing.T) {
	// Test taskQueueManagerControllerData with nil request
	data := taskQueueManagerControllerData{}

	if data.request != nil {
		t.Error("taskQueueManagerControllerData request should be nil")
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
