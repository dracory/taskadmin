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
	"github.com/dracory/taskadmin/task_definition_create"
	"github.com/dracory/taskadmin/task_definition_delete"
	"github.com/dracory/taskadmin/task_definition_manager"
	"github.com/dracory/taskadmin/task_definition_update"
	"github.com/dracory/taskadmin/task_queue_create"
	"github.com/dracory/taskadmin/task_queue_delete"
	"github.com/dracory/taskadmin/task_queue_details"
	"github.com/dracory/taskadmin/task_queue_manager"
	"github.com/dracory/taskadmin/task_queue_parameters"
	"github.com/dracory/taskadmin/task_queue_requeue"
	"github.com/dracory/taskadmin/task_queue_task_restart"
	"github.com/dracory/taskstore"
)

// AdminOptions contains all dependencies and configuration for the task admin.
//
// Store and Logger are required and replace the in-repo version's
// UIOptions fields. Layout is supplied per-request via FuncLayout
// (optional — a default Bootstrap page is used when nil).
type AdminOptions struct {
	// Store is the taskstore.StoreInterface (required)
	Store taskstore.StoreInterface

	// Logger is required
	Logger *slog.Logger

	// FuncLayout returns a fresh shared.Layout for each request.
	// Optional — when nil, a default bare-bones Bootstrap page is used.
	// The returned layout is stateful (SetTitle/SetBody/SetScripts/Render)
	// and must not be shared across requests.
	FuncLayout func(r *http.Request) shared.Layout

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
	funcLayout   func(r *http.Request) shared.Layout
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
		Store:      a.store,
		Logger:     a.logger,
		FuncLayout: a.funcLayout,
	}

	return map[string]func(w http.ResponseWriter, r *http.Request){
		shared.CONTROLLER_HOME:               func(w http.ResponseWriter, r *http.Request) { home.UI(uiConfig).Home(w, r) },
		shared.CONTROLLER_TASK_QUEUE_MANAGER: func(w http.ResponseWriter, r *http.Request) { task_queue_manager.UI(uiConfig).TaskQueueManager(w, r) },
		shared.CONTROLLER_TASK_QUEUE_CREATE:  func(w http.ResponseWriter, r *http.Request) { task_queue_create.UI(uiConfig).TaskQueueCreate(w, r) },
		shared.CONTROLLER_TASK_QUEUE_DELETE:  func(w http.ResponseWriter, r *http.Request) { task_queue_delete.UI(uiConfig).TaskQueueDelete(w, r) },
		shared.CONTROLLER_TASK_QUEUE_DETAILS: func(w http.ResponseWriter, r *http.Request) { task_queue_details.UI(uiConfig).TaskQueueDetails(w, r) },
		shared.CONTROLLER_TASK_QUEUE_PARAMETERS: func(w http.ResponseWriter, r *http.Request) {
			task_queue_parameters.UI(uiConfig).TaskQueueParameters(w, r)
		},
		shared.CONTROLLER_TASK_QUEUE_REQUEUE: func(w http.ResponseWriter, r *http.Request) { task_queue_requeue.UI(uiConfig).TaskQueueRequeue(w, r) },
		shared.CONTROLLER_TASK_QUEUE_TASK_RESTART: func(w http.ResponseWriter, r *http.Request) {
			task_queue_task_restart.UI(uiConfig).TaskQueueTaskRestart(w, r)
		},
		shared.CONTROLLER_TASK_DEFINITION_MANAGER: func(w http.ResponseWriter, r *http.Request) {
			task_definition_manager.UI(uiConfig).TaskDefinitionManager(w, r)
		},
		shared.CONTROLLER_TASK_DEFINITION_CREATE: func(w http.ResponseWriter, r *http.Request) {
			task_definition_create.UI(uiConfig).TaskDefinitionCreate(w, r)
		},
		shared.CONTROLLER_TASK_DEFINITION_UPDATE: func(w http.ResponseWriter, r *http.Request) {
			task_definition_update.UI(uiConfig).TaskDefinitionUpdate(w, r)
		},
		shared.CONTROLLER_TASK_DEFINITION_DELETE: func(w http.ResponseWriter, r *http.Request) {
			task_definition_delete.UI(uiConfig).TaskDefinitionDelete(w, r)
		},
	}
}
