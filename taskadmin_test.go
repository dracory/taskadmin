package taskadmin

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dracory/taskadmin/shared"
	"github.com/dracory/taskadmin/testutils"
)

func Test_New_with_valid_options(t *testing.T) {
	store := testutils.NewStore(t)

	a, err := New(AdminOptions{
		Store:  store,
		Logger: slog.Default(),
	})
	if err != nil {
		t.Fatalf("New() with valid options should not return error, got %v", err)
	}
	if a == nil {
		t.Fatal("New() should return an admin instance")
	}
}

func Test_New_with_missing_options(t *testing.T) {
	if _, err := New(AdminOptions{}); err == nil {
		t.Error("New() with empty options should return error")
	}
}

func Test_New_with_nil_store(t *testing.T) {
	_, err := New(AdminOptions{
		Logger: slog.Default(),
	})
	if err != ErrStoreRequired {
		t.Errorf("New() with nil store should return ErrStoreRequired, got %v", err)
	}
}

func Test_New_with_nil_logger(t *testing.T) {
	store := testutils.NewStore(t)

	_, err := New(AdminOptions{
		Store: store,
	})
	if err != ErrLoggerRequired {
		t.Errorf("New() with nil logger should return ErrLoggerRequired, got %v", err)
	}
}

func Test_Handle_home_controller(t *testing.T) {
	store := testutils.NewStore(t)

	a, err := New(AdminOptions{
		Store:  store,
		Logger: slog.Default(),
	})
	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest("GET", "/admin/task?controller=home", nil)
	w := httptest.NewRecorder()

	a.Handle(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("Handle() should return 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "Dashboard") {
		t.Error("Handle() home should render the dashboard")
	}
}

func Test_Handle_unknown_controller_falls_back_to_home(t *testing.T) {
	store := testutils.NewStore(t)

	a, err := New(AdminOptions{
		Store:  store,
		Logger: slog.Default(),
	})
	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest("GET", "/admin/task?controller=no-such-controller", nil)
	w := httptest.NewRecorder()

	a.Handle(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("Handle() should return 200, got %d", w.Code)
	}
}

func Test_Handle_with_custom_layout(t *testing.T) {
	store := testutils.NewStore(t)
	renderCount := 0

	a, err := New(AdminOptions{
		Store:  store,
		Logger: slog.Default(),
		FuncLayout: func(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options shared.LayoutOptions) string {
			renderCount++
			return "<html><head><title>" + webpageTitle + "</title></head><body>" + webpageHtml + "</body></html>"
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest("GET", "/admin/task?controller=home", nil)
	w := httptest.NewRecorder()

	a.Handle(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("Handle() should return 200, got %d", w.Code)
	}
	if renderCount == 0 {
		t.Error("Handle() should use the custom layout")
	}
}

func Test_Handle_auth_redirect(t *testing.T) {
	store := testutils.NewStore(t)

	a, err := New(AdminOptions{
		Store:        store,
		Logger:       slog.Default(),
		AdminHomeURL: "/admin",
		AuthUserID: func(r *http.Request) string {
			return ""
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest("GET", "/admin/task?controller=home", nil)
	w := httptest.NewRecorder()

	a.Handle(w, r)

	if w.Code != http.StatusSeeOther {
		t.Errorf("Handle() should redirect unauthenticated users, got %d", w.Code)
	}
	if w.Header().Get("Location") != "/admin" {
		t.Errorf("Handle() should redirect to /admin, got %s", w.Header().Get("Location"))
	}
}

func Test_AdminOptions_defaults(t *testing.T) {
	store := testutils.NewStore(t)

	a, err := New(AdminOptions{
		Store:  store,
		Logger: slog.Default(),
	})
	if err != nil {
		t.Fatal(err)
	}

	admin := a.(*admin)
	if admin.adminHomeURL != "/admin" {
		t.Errorf("adminHomeURL should default to /admin, got %s", admin.adminHomeURL)
	}
	if admin.taskAdminURL != "/admin/task" {
		t.Errorf("taskAdminURL should default to /admin/task, got %s", admin.taskAdminURL)
	}
}
