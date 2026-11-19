package home_test

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/dracory/taskadmin/home"
	"github.com/dracory/taskadmin/shared"
	"github.com/dracory/taskadmin/testutils"
)

func Test_UI(t *testing.T) {
	ui := home.UI(shared.UiConfig{
		Store:  testutils.NewStore(t),
		Logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
	})
	if ui == nil {
		t.Fatal("UI() should return a non-nil controller")
	}
	if ui.Store() == nil {
		t.Fatal("UI() should set the store")
	}
}

func Test_renderPage(t *testing.T) {
	ui := home.UI(shared.UiConfig{
		Store:  testutils.NewStore(t),
		Logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
		Layout: func(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options shared.LayoutOptions) string {
			return "<html><head><title>" + webpageTitle + "</title></head><body>" + webpageHtml + "</body></html>"
		},
	})

	r := httptest.NewRequest("GET", "/admin/task?controller=home", nil)
	w := httptest.NewRecorder()
	ui.Home(w, r)

	if w.Code != 200 {
		t.Fatalf("page should return 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "Tasks. Dashboard") {
		t.Error("page should render the dashboard header")
	}
}
