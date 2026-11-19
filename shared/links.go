package shared

import (
	"net/http"
)

// Links provides URL helpers for taskadmin controllers.
// The base URL is read from request context (injected by Handle()),
// not hardcoded.
type Links struct {
	baseURL string
}

// NewLinks creates a Links helper with the given base URL.
// If baseURL is empty, defaults to "/admin/task".
func NewLinks(baseURL string) *Links {
	if baseURL == "" {
		baseURL = "/admin/task"
	}
	return &Links{baseURL: baseURL}
}

// NewLinksFromRequest creates a Links helper using the task admin URL
// from the request context.
func NewLinksFromRequest(r *http.Request) *Links {
	return NewLinks(TaskAdminURL(r))
}

// Home builds the URL for the home controller
func (l *Links) Home(params map[string]string) string {
	return l.url(CONTROLLER_HOME, params)
}

// TaskQueue builds the URL for the task queue manager controller
func (l *Links) TaskQueue(params map[string]string) string {
	return l.url(CONTROLLER_TASK_QUEUE_MANAGER, params)
}

// TaskDefinitions builds the URL for the task definition manager controller
func (l *Links) TaskDefinitions(params map[string]string) string {
	return l.url(CONTROLLER_TASK_DEFINITION_MANAGER, params)
}

// url builds a controller URL on the task admin endpoint
func (l *Links) url(controller string, params map[string]string) string {
	if params == nil {
		params = map[string]string{}
	}
	params["controller"] = controller
	return l.baseURL + Query(params)
}
