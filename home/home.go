package home

import (
	"net/http"

	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/taskadmin/shared"
)

// UiInterface defines the home controller's UI interface
type UiInterface interface {
	shared.UiInterface
	Home(w http.ResponseWriter, r *http.Request)
}

// ui implements UiInterface
type ui struct {
	shared.UiBase
}

// UI creates a new home controller UI from the given config
func UI(config shared.UiConfig) UiInterface {
	return &ui{UiBase: shared.NewUiBase(config)}
}

// Home handles the home controller requests
func (u *ui) Home(w http.ResponseWriter, r *http.Request) {
	html := u.Handler(w, r)
	if html != "" {
		if w.Header().Get("Content-Type") == "" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
		}
		// #nosec G705 -- html is rendered by the page builder, not user input
		_, _ = w.Write([]byte(html))
	}
}

// Handler processes the home controller request and returns HTML
func (u *ui) Handler(w http.ResponseWriter, r *http.Request) string {
	return u.renderPage(w, r)
}

func (u *ui) renderPage(w http.ResponseWriter, r *http.Request) string {
	if u.Store() == nil {
		return shared.ErrorAlert("Task store is not initialized")
	}

	breadcrumbs := shared.Breadcrumbs(r, []shared.Breadcrumb{
		{Name: "Dashboard", URL: ""},
	})

	content := hb.Div().
		Class("container").
		Class("py-4").
		Child(breadcrumbs).
		Child(shared.PageHeader("Tasks. Dashboard")).
		Child(shared.SectionNav(r, u.Store(), u.Logger()))

	return u.Layout()(w, r, "Dashboard | Tasks", content.ToHTML(), shared.LayoutOptions{
		ScriptURLs: []string{
			cdn.VueJs_3(),
		},
	})
}
