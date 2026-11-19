package task_definition_delete

import (
	"github.com/dracory/taskadmin/shared"
	"github.com/dracory/taskadmin/testutils"
	"github.com/dracory/taskstore"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_taskDefinitionDelete(t *testing.T) {
	// Test taskDefinitionDelete controller constructor with real SQLite store
	store := testutils.NewStore(t)

	controller := testUI(store)

	if controller == nil {
		t.Error("taskDefinitionDelete() should return a non-nil controller")
	}
	if controller.Store() == nil {
		t.Error("taskDefinitionDelete() should set store")
	}
}

func Test_taskDefinitionDeleteController_struct_fields(t *testing.T) {
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

func Test_taskDefinitionDeleteController_with_nil_store(t *testing.T) {
	// Test ui with nil store

	controller := testUI(nil)

	if controller.Store() != nil {
		t.Error("ui store should be nil")
	}
}

func Test_taskDefinitionDeleteControllerData_request_field(t *testing.T) {
	// Test taskDefinitionDeleteControllerData request field
	req := httptest.NewRequest("GET", "/", nil)
	data := taskDefinitionDeleteControllerData{request: req}

	if data.request == nil {
		t.Error("taskDefinitionDeleteControllerData request should not be nil")
	}
}

func Test_taskDefinitionDeleteControllerData_nil_request(t *testing.T) {
	// Test taskDefinitionDeleteControllerData with nil request
	data := taskDefinitionDeleteControllerData{}

	if data.request != nil {
		t.Error("taskDefinitionDeleteControllerData request should be nil")
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
