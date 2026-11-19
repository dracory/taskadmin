package home

import (
	"context"
	"github.com/dracory/taskadmin/shared"
	"github.com/dracory/taskadmin/testutils"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dracory/taskstore"
)

func Test_home(t *testing.T) {
	// Test home controller constructor with real SQLite store
	store := testutils.NewStore(t)

	controller := testUI(store)

	if controller == nil {
		t.Error("home() should return a non-nil controller")
	}
	if controller.Store() == nil {
		t.Error("home() should set store")
	}
	if controller.Layout(nil) == nil {
		t.Error("home() should set layout")
	}
}

func Test_homeController_prepareData(t *testing.T) {
	// Test prepareData method
	store := testutils.NewStore(t)

	controller := testUI(store)

	req := httptest.NewRequest("GET", "/", nil)
	data, errorMessage := controller.prepareData(req)

	if data.request == nil {
		t.Error("prepareData() should set request")
	}
	if errorMessage != "" {
		t.Error("prepareData() should not return error message")
	}
}

func Test_homeController_prepareData_with_nil_request(t *testing.T) {
	// Test prepareData with nil request
	store := testutils.NewStore(t)

	controller := testUI(store)

	data, errorMessage := controller.prepareData(nil)

	if data.request != nil {
		t.Error("prepareData() with nil request should have nil request in data")
	}
	if errorMessage != "" {
		t.Error("prepareData() with nil request should not return error message")
	}
}

func Test_homeControllerData_request_field(t *testing.T) {
	// Test homeControllerData request field
	req := httptest.NewRequest("GET", "/", nil)
	data := homeControllerData{request: req}

	if data.request == nil {
		t.Error("homeControllerData request should not be nil")
	}
}

func Test_homeControllerData_nil_request(t *testing.T) {
	// Test homeControllerData with nil request
	data := homeControllerData{}

	if data.request != nil {
		t.Error("homeControllerData request should be nil")
	}
}

func Test_homeController_with_real_store(t *testing.T) {
	// Test homeController with real SQLite store
	store := testutils.NewStore(t)

	controller := testUI(store)

	// Create a task definition in the store
	ctx := context.Background()
	taskDef := taskstore.NewTaskDefinition()
	taskDef.SetTitle("Test Task")
	taskDef.SetAlias("test-task")
	taskDef.SetDescription("Test description")

	err := store.TaskDefinitionCreate(ctx, taskDef)
	if err != nil {
		t.Fatalf("TaskDefinitionCreate() failed: %v", err)
	}

	// Verify controller has access to store
	if controller.Store() == nil {
		t.Error("homeController store should not be nil")
	}
}

func Test_homeController_with_nil_store(t *testing.T) {
	// Test homeController with nil store

	controller := testUI(nil)

	if controller.Store() != nil {
		t.Error("homeController store should be nil")
	}
	if controller.Layout(nil) == nil {
		t.Error("homeController layout should not be nil")
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
