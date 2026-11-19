package task_queue_details

import (
	"github.com/dracory/taskadmin/shared"
	"github.com/dracory/taskadmin/testutils"
	"github.com/dracory/taskstore"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_taskQueueDetails(t *testing.T) {
	// Test taskQueueDetails controller constructor with real SQLite store
	store := testutils.NewStore(t)

	controller := testUI(store)

	if controller == nil {
		t.Error("taskQueueDetails() should return a non-nil controller")
	}
	if controller.Store() == nil {
		t.Error("taskQueueDetails() should set store")
	}
}

func Test_taskQueueDetailsController_struct_fields(t *testing.T) {
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

func Test_taskQueueDetailsController_with_nil_store(t *testing.T) {
	// Test ui with nil store

	controller := testUI(nil)

	if controller.Store() != nil {
		t.Error("ui store should be nil")
	}
}

func Test_taskQueueDetailsControllerData_request_field(t *testing.T) {
	// Test taskQueueDetailsControllerData request field
	req := httptest.NewRequest("GET", "/", nil)
	data := taskQueueDetailsControllerData{request: req}

	if data.request == nil {
		t.Error("taskQueueDetailsControllerData request should not be nil")
	}
}

func Test_taskQueueDetailsControllerData_nil_request(t *testing.T) {
	// Test taskQueueDetailsControllerData with nil request
	data := taskQueueDetailsControllerData{}

	if data.request != nil {
		t.Error("taskQueueDetailsControllerData request should be nil")
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
