package schedule_manager

import (
	"net/http"

	"github.com/dracory/req"
	"github.com/dracory/taskadmin/shared"
)

const (
	actionLoadSchedules   = "load-schedules"
	actionLoadDefinitions = "load-definitions"
	actionScheduleCreate  = "create-schedule"
	actionScheduleUpdate  = "update-schedule"
	actionScheduleDelete  = "delete-schedule"
	actionScheduleDetails = "schedule-details"
)

// UiInterface defines the schedule manager controller's UI interface
type UiInterface interface {
	shared.UiInterface
	ScheduleManager(w http.ResponseWriter, r *http.Request)
}

// ui implements UiInterface
type ui struct {
	shared.UiBase
}

// UI creates a new schedule manager controller UI from the given config
func UI(config shared.UiConfig) UiInterface {
	return &ui{UiBase: shared.NewUiBase(config)}
}

// ScheduleManager handles the schedule manager controller requests
func (u *ui) ScheduleManager(w http.ResponseWriter, r *http.Request) {
	html := u.Handler(w, r)
	if html != "" {
		if w.Header().Get("Content-Type") == "" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
		}
		// #nosec G705 -- html is rendered by the page builder, not user input
		_, _ = w.Write([]byte(html))
	}
}

// Handler processes the schedule manager request and returns HTML
func (u *ui) Handler(w http.ResponseWriter, r *http.Request) string {
	action := req.GetStringTrimmed(r, "action")

	switch action {
	case actionLoadSchedules:
		return u.handleLoadSchedules(w, r)
	case actionLoadDefinitions:
		return u.handleLoadDefinitions(w, r)
	case actionScheduleCreate:
		return u.handleScheduleCreate(w, r)
	case actionScheduleUpdate:
		return u.handleScheduleUpdate(w, r)
	case actionScheduleDelete:
		return u.handleScheduleDelete(w, r)
	case actionScheduleDetails:
		return u.handleScheduleDetails(w, r)
	default:
		return u.renderPage(w, r)
	}
}
