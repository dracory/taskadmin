package shared

import "net/http"

// AdminHomeURL returns the admin home URL from context
func AdminHomeURL(r *http.Request) string {
	value, ok := r.Context().Value(KeyAdminHomeURL).(string)
	if !ok || value == "" {
		return "/admin"
	}
	return value
}

// TaskAdminURL returns the task admin base URL from context
func TaskAdminURL(r *http.Request) string {
	value, ok := r.Context().Value(KeyTaskAdminURL).(string)
	if !ok || value == "" {
		return "/admin/task"
	}
	return value
}
