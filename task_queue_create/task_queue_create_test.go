package task_queue_create

import (
	"github.com/dracory/taskadmin/shared"
	"github.com/dracory/taskadmin/testutils"
	"github.com/dracory/taskstore"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_taskQueueCreate(t *testing.T) {
	// Test taskQueueCreate controller constructor with real SQLite store
	store := testutils.NewStore(t)

	controller := testUI(store)

	if controller == nil {
		t.Error("taskQueueCreate() should return a non-nil controller")
	}
	if controller.Store() == nil {
		t.Error("taskQueueCreate() should set store")
	}
}

func Test_taskQueueCreateController_struct_fields(t *testing.T) {
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

func Test_taskQueueCreateController_with_nil_store(t *testing.T) {
	// Test ui with nil store

	controller := testUI(nil)

	if controller.Store() != nil {
		t.Error("ui store should be nil")
	}
}

func Test_taskQueueCreateControllerData_request_field(t *testing.T) {
	// Test taskQueueCreateControllerData request field
	req := httptest.NewRequest("GET", "/", nil)
	data := taskQueueCreateControllerData{request: req}

	if data.request == nil {
		t.Error("taskQueueCreateControllerData request should not be nil")
	}
}

func Test_taskQueueCreateControllerData_nil_request(t *testing.T) {
	// Test taskQueueCreateControllerData with nil request
	data := taskQueueCreateControllerData{}

	if data.request != nil {
		t.Error("taskQueueCreateControllerData request should be nil")
	}
}

func Test_taskQueueCreateController_with_real_store(t *testing.T) {
	// Test ui with real SQLite store
	store := testutils.NewStore(t)

	controller := testUI(store)

	// Verify controller has access to store
	if controller.Store() == nil {
		t.Error("ui store should not be nil")
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
