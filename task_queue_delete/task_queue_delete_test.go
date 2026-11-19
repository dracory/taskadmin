package task_queue_delete

import (
	"github.com/dracory/taskadmin/shared"
	"github.com/dracory/taskadmin/testutils"
	"github.com/dracory/taskstore"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_taskQueueDelete(t *testing.T) {
	// Test taskQueueDelete controller constructor with real SQLite store
	store := testutils.NewStore(t)

	controller := testUI(store)

	if controller == nil {
		t.Error("taskQueueDelete() should return a non-nil controller")
	}
	if controller.Store() == nil {
		t.Error("taskQueueDelete() should set store")
	}
}

func Test_taskQueueDeleteController_struct_fields(t *testing.T) {
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

func Test_taskQueueDeleteController_with_nil_store(t *testing.T) {
	// Test ui with nil store

	controller := testUI(nil)

	if controller.Store() != nil {
		t.Error("ui store should be nil")
	}
}

func Test_taskQueueDeleteControllerData_request_field(t *testing.T) {
	// Test taskQueueDeleteControllerData request field
	req := httptest.NewRequest("GET", "/", nil)
	data := taskQueueDeleteControllerData{request: req}

	if data.request == nil {
		t.Error("taskQueueDeleteControllerData request should not be nil")
	}
}

func Test_taskQueueDeleteControllerData_nil_request(t *testing.T) {
	// Test taskQueueDeleteControllerData with nil request
	data := taskQueueDeleteControllerData{}

	if data.request != nil {
		t.Error("taskQueueDeleteControllerData request should be nil")
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
