package home

import (
	"embed"
	"net/http"

	"github.com/dracory/api"
	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/taskadmin/shared"
	"github.com/dracory/taskstore"
)

//go:embed *.html
//go:embed *.js
var homeFiles embed.FS

const (
	actionLoadStats = "load-stats"
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
	action := r.URL.Query().Get("action")

	switch action {
	case actionLoadStats:
		return u.handleLoadStats(w, r)
	default:
		return u.renderPage(w, r)
	}
}

func (u *ui) renderPage(w http.ResponseWriter, r *http.Request) string {
	if u.Store() == nil {
		return shared.ErrorAlert("Task store is not initialized")
	}

	breadcrumbs := shared.Breadcrumbs(r, []shared.Breadcrumb{
		{Name: "Dashboard", URL: ""},
	})

	htmlContent, err := homeFiles.ReadFile("home.html")
	if err != nil {
		u.Logger().Error("Failed to read home HTML template", "error", err)
		return shared.ErrorAlert("Error loading home component")
	}

	jsContent, err := homeFiles.ReadFile("home.js")
	if err != nil {
		u.Logger().Error("Failed to read home JavaScript file", "error", err)
		return shared.ErrorAlert("Error loading home component")
	}

	initScript := hb.Script(`
		const urlLoadStats = '` + shared.URL(r, shared.CONTROLLER_HOME, map[string]string{"action": actionLoadStats}) + `';
		const urlTaskQueue = '` + shared.URL(r, shared.CONTROLLER_TASK_QUEUE_MANAGER, nil) + `';
		const urlTaskDefinitions = '` + shared.URL(r, shared.CONTROLLER_TASK_DEFINITION_MANAGER, nil) + `';
		const urlSchedules = '` + shared.URL(r, shared.CONTROLLER_SCHEDULE_MANAGER, nil) + `';
	`)

	vueContainer := hb.Div().
		Child(hb.Script("").Src(cdn.VueJs_3())).
		Child(hb.Wrap().HTML(string(htmlContent))).
		Child(initScript).
		Child(hb.Script(string(jsContent)))

	content := hb.Div().
		Class("container").
		Class("py-4").
		Child(shared.SectionNav(r, u.Store(), u.Logger())).
		Child(breadcrumbs).
		Child(shared.PageHeader("Tasks. Dashboard")).
		Child(vueContainer)

	return u.Layout()(w, r, "Dashboard | Tasks", content.ToHTML(), shared.LayoutOptions{})
}

// handleLoadStats returns queue and definition counts as JSON
func (u *ui) handleLoadStats(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	countFor := func(query taskstore.TaskQueueQueryInterface) int64 {
		count, err := u.Store().TaskQueueCount(ctx, query)
		if err != nil {
			u.Logger().Error("Failed to count queued tasks", "error", err)
			return 0
		}
		return count
	}

	queuedCount := countFor(taskstore.TaskQueueQuery().SetStatus(taskstore.TaskQueueStatusQueued))
	runningCount := countFor(taskstore.TaskQueueQuery().SetStatus(taskstore.TaskQueueStatusRunning))
	failedCount := countFor(taskstore.TaskQueueQuery().SetStatus(taskstore.TaskQueueStatusFailed))

	definitionCount, err := u.Store().TaskDefinitionCount(ctx, taskstore.TaskDefinitionQuery())
	if err != nil {
		u.Logger().Error("Failed to count task definitions", "error", err)
		definitionCount = 0
	}

	scheduleCount, err := u.Store().ScheduleCount(ctx, taskstore.NewScheduleQuery().SetLimit(0))
	if err != nil {
		u.Logger().Error("Failed to count schedules", "error", err)
		scheduleCount = 0
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(api.SuccessWithData("Stats loaded successfully", map[string]any{
		"queued_count":     queuedCount,
		"running_count":    runningCount,
		"failed_count":     failedCount,
		"definition_count": definitionCount,
		"schedule_count":   scheduleCount,
	}).ToString()))
	return ""
}
