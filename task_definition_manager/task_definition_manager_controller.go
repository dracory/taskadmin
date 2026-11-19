package task_definition_manager

import (
	"net/http"

	"github.com/dracory/req"
	"github.com/dracory/taskadmin/shared"
)

const (
	actionLoadDefinitions   = "load-definitions"
	actionDefinitionCreate  = "create-definition"
	actionDefinitionUpdate  = "update-definition"
	actionDefinitionDelete  = "delete-definition"
	actionDefinitionDetails = "definition-details"
)

// UiInterface defines the task definition manager controller's UI interface
type UiInterface interface {
	shared.UiInterface
	TaskDefinitionManager(w http.ResponseWriter, r *http.Request)
}

// ui implements UiInterface
type ui struct {
	shared.UiBase
}

// UI creates a new task definition manager controller UI from the given config
func UI(config shared.UiConfig) UiInterface {
	return &ui{UiBase: shared.NewUiBase(config)}
}

// TaskDefinitionManager handles the task definition manager controller requests
func (u *ui) TaskDefinitionManager(w http.ResponseWriter, r *http.Request) {
	html := u.Handler(w, r)
	if html != "" {
		if w.Header().Get("Content-Type") == "" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
		}
		// #nosec G705 -- html is rendered by the page builder, not user input
		_, _ = w.Write([]byte(html))
	}
}

// Handler processes the task definition manager request and returns HTML
func (u *ui) Handler(w http.ResponseWriter, r *http.Request) string {
	action := req.GetStringTrimmed(r, "action")

	switch action {
	case actionLoadDefinitions:
		return u.handleLoadDefinitions(w, r)
	case actionDefinitionCreate:
		return u.handleDefinitionCreate(w, r)
	case actionDefinitionUpdate:
		return u.handleDefinitionUpdate(w, r)
	case actionDefinitionDelete:
		return u.handleDefinitionDelete(w, r)
	case actionDefinitionDetails:
		return u.handleDefinitionDetails(w, r)
	default:
		return u.renderPage(w, r)
	}
}
