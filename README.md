# taskadmin

Standalone admin interface for [taskstore](https://github.com/dracory/taskstore),
following the folder-per-controller pattern used by `shopadmin`.

Extracted from `github.com/dracory/taskstore/admin`.

## Structure

- `taskadmin.go` — `AdminOptions`, `New()`, `Handle()` dispatcher
- `shared/` — `UiConfig`/`UiBase`, `FuncLayout` + default layout,
  breadcrumbs, page header, section nav, URL/query helpers
- `home/` — dashboard landing page
- `task_queue_manager/` — queued task manager. Vue page (`tasks.html` +
  `tasks.js`) plus JSON AJAX handlers (`handle_*_ajax.go`) for loading,
  enqueueing, restarting, requeueing, deleting, and editing task parameters
- `task_definition_manager/` — task definition manager. Vue page
  (`definitions.html` + `definitions.js`) plus JSON AJAX handlers for
  loading, creating, updating, and deleting definitions
- `testutils/` — in-memory SQLite store for tests

Each manager is a thin `?action=` dispatcher: with no action it renders the
page (layout + embedded Vue assets); with an action it returns a JSON
response via `github.com/dracory/api`.

## Usage

```go
admin, err := taskadmin.New(taskadmin.AdminOptions{
    Store:  store,           // taskstore.StoreInterface (required)
    Logger: slog.Default(),  // required

    // optional
    FuncLayout: func(w http.ResponseWriter, r *http.Request,
        webpageTitle, webpageHtml string, options shared.LayoutOptions) string {
        return myAdminLayout(w, r, webpageTitle, webpageHtml, options)
    },
    AdminHomeURL: "/admin",      // default
    TaskAdminURL: "/admin/task", // default
    AuthUserID: func(r *http.Request) string {
        return authUserID(r) // "" => redirect to AdminHomeURL
    },
})
if err != nil {
    panic(err)
}

mux.HandleFunc("/admin/task", admin.Handle)
```

Routing is by `?controller=` query parameter (`home`, `task-queue-manager`,
`task-definition-manager`), defaulting to `home`.

## Testing

```bash
go test ./...
# or
task test
```
