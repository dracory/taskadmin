// Package taskadmin provides a standalone task admin interface following
// the folder-per-controller pattern. Each controller is in its own
// subfolder and handles its own views and AJAX data.
//
// This module is modeled on github.com/dracory/shopadmin and was
// extracted from github.com/dracory/taskstore/admin.
package taskadmin

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/dracory/req"
	"github.com/dracory/taskadmin/home"
	"github.com/dracory/taskadmin/shared"
	"github.com/dracory/taskadmin/task_definition_manager"
	"github.com/dracory/taskadmin/task_queue_manager"
	"github.com/dracory/taskstore"
)

// AdminOptions contains all dependencies and configuration for the task admin.
type AdminOptions struct {
	// Store is the taskstore.StoreInterface (required)
	Store taskstore.StoreInterface

	// Logger is required
	Logger *slog.Logger

	// FuncLayout is an optional function to render the admin interface
	// inside your own layout (branding, menus, etc.). It receives the
	// request and response writer so the host project can access request
	// context (auth user, locale, etc.) when rendering the layout.
	// If nil, a default bare-bones HTML page is used (Bootstrap CDN).
	FuncLayout shared.FuncLayout

	// AdminHomeURL is the URL for the admin home page (default: "/admin")
	AdminHomeURL string

	// TaskAdminURL is the base URL for the task admin (default: "/admin/task")
	TaskAdminURL string

	// AuthUserID returns the authenticated user ID from the request.
	// If it returns "", the user is treated as unauthenticated.
	// Optional — nil disables the auth check.
	AuthUserID func(r *http.Request) string
}

// AdminInterface defines the interface for the task admin
type AdminInterface interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

// admin implements AdminInterface
type admin struct {
	store        taskstore.StoreInterface
	logger       *slog.Logger
	funcLayout   shared.FuncLayout
	adminHomeURL string
	taskAdminURL string
	authUserID   func(r *http.Request) string
	routes       map[string]func(w http.ResponseWriter, r *http.Request)
}

// New creates a new task admin instance.
// Returns ErrStoreRequired if Store is nil, ErrLoggerRequired if Logger is nil.
func New(opts AdminOptions) (AdminInterface, error) {
	if opts.Store == nil {
		return nil, ErrStoreRequired
	}
	if opts.Logger == nil {
		return nil, ErrLoggerRequired
	}

	if opts.AdminHomeURL == "" {
		opts.AdminHomeURL = "/admin"
	}
	if opts.TaskAdminURL == "" {
		opts.TaskAdminURL = "/admin/task"
	}

	a := &admin{
		store:        opts.Store,
		logger:       opts.Logger,
		funcLayout:   opts.FuncLayout,
		adminHomeURL: opts.AdminHomeURL,
		taskAdminURL: opts.TaskAdminURL,
		authUserID:   opts.AuthUserID,
	}

	a.routes = a.buildRoutes()

	return a, nil
}

// Handle processes all task admin requests.
// Config values are injected into the request context (following the
// shopadmin pattern). Route lookup is map-based on the "controller"
// query parameter.
func (a *admin) Handle(w http.ResponseWriter, r *http.Request) {
	if a.authUserID != nil && a.authUserID(r) == "" {
		http.Redirect(w, r, a.adminHomeURL, http.StatusSeeOther)
		return
	}

	ctx := context.WithValue(r.Context(), shared.KeyEndpoint, r.URL.Path)
	ctx = context.WithValue(ctx, shared.KeyAdminHomeURL, a.adminHomeURL)
	ctx = context.WithValue(ctx, shared.KeyTaskAdminURL, a.taskAdminURL)
	r = r.WithContext(ctx)

	controller := req.GetStringTrimmed(r, "controller")
	if controller == "" {
		controller = shared.CONTROLLER_HOME
	}

	handler, ok := a.routes[controller]
	if !ok {
		handler = a.routes[shared.CONTROLLER_HOME]
	}

	handler(w, r)
}

// buildRoutes creates the handler dispatch map once at construction time.
func (a *admin) buildRoutes() map[string]func(w http.ResponseWriter, r *http.Request) {
	uiConfig := shared.UiConfig{
		Store:  a.store,
		Logger: a.logger,
		Layout: a.render,
	}

	return map[string]func(w http.ResponseWriter, r *http.Request){
		shared.CONTROLLER_HOME:               func(w http.ResponseWriter, r *http.Request) { home.UI(uiConfig).Home(w, r) },
		shared.CONTROLLER_TASK_QUEUE_MANAGER: func(w http.ResponseWriter, r *http.Request) { task_queue_manager.UI(uiConfig).TaskQueueManager(w, r) },
		shared.CONTROLLER_TASK_DEFINITION_MANAGER: func(w http.ResponseWriter, r *http.Request) {
			task_definition_manager.UI(uiConfig).TaskDefinitionManager(w, r)
		},
	}
}

// render wraps content in the layout. If FuncLayout is provided and
// returns non-empty, it is used; otherwise the default shared.Layout
// is used (following the shopadmin pattern).
//
// When FuncLayout is set, the default shared.Layout is NOT computed.
func (a *admin) render(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options shared.LayoutOptions) string {
	if a.funcLayout != nil {
		custom := a.funcLayout(w, r, webpageTitle, webpageHtml, options)
		if custom != "" {
			if w != nil {
				if w.Header().Get("Content-Type") == "" {
					w.Header().Set("Content-Type", "text/html; charset=utf-8")
				}
				_, _ = w.Write([]byte(custom))
				return ""
			}
			return custom
		}
	}

	webpage := shared.Layout(w, r, webpageTitle, webpageHtml, options)

	if w != nil {
		if w.Header().Get("Content-Type") == "" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
		}
		_, _ = w.Write([]byte(webpage))
		return ""
	}

	return webpage
}
