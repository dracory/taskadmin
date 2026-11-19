package task_queue_manager

import (
	"net/http"

	"github.com/dracory/req"
	"github.com/dracory/taskadmin/shared"
)

const (
	actionLoadTasks          = "load-tasks"
	actionTaskEnqueue        = "enqueue-task"
	actionTaskDelete         = "delete-task"
	actionTaskDeleteSelected = "delete-selected"
	actionTaskRequeue        = "requeue-task"
	actionTaskRestart        = "restart-task"
	actionTaskDetails        = "task-details"
	actionTaskUpdateParams   = "update-parameters"
)

// UiInterface defines the task queue manager controller's UI interface
type UiInterface interface {
	shared.UiInterface
	TaskQueueManager(w http.ResponseWriter, r *http.Request)
}

// ui implements UiInterface
type ui struct {
	shared.UiBase
}

// UI creates a new task queue manager controller UI from the given config
func UI(config shared.UiConfig) UiInterface {
	return &ui{UiBase: shared.NewUiBase(config)}
}

// TaskQueueManager handles the task queue manager controller requests
func (u *ui) TaskQueueManager(w http.ResponseWriter, r *http.Request) {
	html := u.Handler(w, r)
	if html != "" {
		if w.Header().Get("Content-Type") == "" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
		}
		// #nosec G705 -- html is rendered by the page builder, not user input
		_, _ = w.Write([]byte(html))
	}
}

// Handler processes the task queue manager request and returns HTML
func (u *ui) Handler(w http.ResponseWriter, r *http.Request) string {
	action := req.GetStringTrimmed(r, "action")

	switch action {
	case actionLoadTasks:
		return u.handleLoadTasks(w, r)
	case actionTaskEnqueue:
		return u.handleTaskEnqueue(w, r)
	case actionTaskDelete:
		return u.handleTaskDelete(w, r)
	case actionTaskDeleteSelected:
		return u.handleTaskDeleteSelected(w, r)
	case actionTaskRequeue:
		return u.handleTaskRequeue(w, r)
	case actionTaskRestart:
		return u.handleTaskRestart(w, r)
	case actionTaskDetails:
		return u.handleTaskDetails(w, r)
	case actionTaskUpdateParams:
		return u.handleTaskUpdateParameters(w, r)
	default:
		return u.renderPage(w, r)
	}
}
