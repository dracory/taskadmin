package task_definition_update

import (
	"github.com/dracory/taskadmin/shared"
	"github.com/dracory/taskadmin/testutils"
	"github.com/dracory/taskstore"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_taskDefinitionUpdate(t *testing.T) {
	// Test taskDefinitionUpdate controller constructor with real SQLite store
	store := testutils.NewStore(t)

	controller := testUI(store)

	if controller == nil {
		t.Error("taskDefinitionUpdate() should return a non-nil controller")
	}
	if controller.Store() == nil {
		t.Error("taskDefinitionUpdate() should set store")
	}
}

func Test_taskDefinitionUpdateController_struct_fields(t *testing.T) {
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

func Test_taskDefinitionUpdateController_with_nil_store(t *testing.T) {
	// Test ui with nil store

	controller := testUI(nil)

	if controller.Store() != nil {
		t.Error("ui store should be nil")
	}
}

func Test_taskDefinitionUpdateControllerData_request_field(t *testing.T) {
	// Test taskDefinitionUpdateControllerData request field
	req := httptest.NewRequest("GET", "/", nil)
	data := taskDefinitionUpdateControllerData{request: req}

	if data.request == nil {
		t.Error("taskDefinitionUpdateControllerData request should not be nil")
	}
}

func Test_taskDefinitionUpdateControllerData_nil_request(t *testing.T) {
	// Test taskDefinitionUpdateControllerData with nil request
	data := taskDefinitionUpdateControllerData{}

	if data.request != nil {
		t.Error("taskDefinitionUpdateControllerData request should be nil")
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
