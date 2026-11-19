package task_queue_requeue

import (
	"github.com/dracory/taskadmin/shared"
	"github.com/dracory/taskadmin/testutils"
	"github.com/dracory/taskstore"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_taskQueueRequeue(t *testing.T) {
	// Test taskQueueRequeue controller constructor with real SQLite store
	store := testutils.NewStore(t)

	controller := testUI(store)

	if controller == nil {
		t.Error("taskQueueRequeue() should return a non-nil controller")
	}
	if controller.Store() == nil {
		t.Error("taskQueueRequeue() should set store")
	}
}

func Test_taskQueueRequeueController_struct_fields(t *testing.T) {
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

func Test_taskQueueRequeueController_with_nil_store(t *testing.T) {
	// Test ui with nil store

	controller := testUI(nil)

	if controller.Store() != nil {
		t.Error("ui store should be nil")
	}
}

func Test_taskQueueRequeueControllerData_request_field(t *testing.T) {
	// Test taskQueueRequeueControllerData request field
	req := httptest.NewRequest("GET", "/", nil)
	data := taskQueueRequeueControllerData{request: req}

	if data.request == nil {
		t.Error("taskQueueRequeueControllerData request should not be nil")
	}
}

func Test_taskQueueRequeueControllerData_nil_request(t *testing.T) {
	// Test taskQueueRequeueControllerData with nil request
	data := taskQueueRequeueControllerData{}

	if data.request != nil {
		t.Error("taskQueueRequeueControllerData request should be nil")
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
