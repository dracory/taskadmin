# taskadmin

Standalone admin interface for [taskstore](https://github.com/dracory/taskstore),
following the folder-per-controller pattern used by `shopadmin` and `cmsstore/admin`.

Extracted from `github.com/dracory/taskstore/admin`.

## Structure

- `taskadmin.go` — `AdminOptions`, `New()`, `Handle()` dispatcher
- `shared/` — UiConfig/UiBase, Layout interface + default webpage,
  breadcrumbs, page header, nav tabs, URL/query helpers
- `home/` — dashboard landing page
- `task_queue_manager/` — queued task list (filters, sorting, actions)
- `task_queue_create/`, `task_queue_delete/`, `task_queue_details/`,
  `task_queue_parameters/`, `task_queue_requeue/`,
  `task_queue_task_restart/` — queued task modals/actions
- `task_definition_manager/` — task definition list
- `task_definition_create/`, `task_definition_update/`,
  `task_definition_delete/` — task definition modals/actions
- `testutils/` — in-memory SQLite store + mock layout for tests

## Usage

```go
admin, err := taskadmin.New(taskadmin.AdminOptions{
    Store:  store,           // taskstore.StoreInterface (required)
    Logger: slog.Default(),  // required

    // optional
    FuncLayout: func(r *http.Request) shared.Layout {
        return myAdminLayout // fresh instance per request
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

Routing is by `?controller=` query parameter
(`home`, `task-queue-manager`, `task-definition-manager`, and the
modal/action controllers listed above).

## Testing

```bash
go test ./...
# or
task test
```
