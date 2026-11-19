package shared

import (
	"log/slog"
	"net/http"

	"github.com/dracory/hb"
	"github.com/dracory/taskstore"
	"github.com/spf13/cast"
)

// SectionNav renders the task admin section navigation as pills,
// matching the pattern used across the admin packages. The active tab
// is derived from the ?controller= query parameter. The Task Queue and
// Task Definitions tabs show record counts as badges.
func SectionNav(r *http.Request, store taskstore.StoreInterface, logger *slog.Logger) hb.TagInterface {
	linksHelper := NewLinksFromRequest(r)

	active := "dashboard"
	controller := ""
	if r != nil && r.URL != nil {
		controller = r.URL.Query().Get("controller")
	}
	switch controller {
	case CONTROLLER_TASK_QUEUE_MANAGER:
		active = "queue"
	case CONTROLLER_TASK_DEFINITION_MANAGER:
		active = "definitions"
	case CONTROLLER_SCHEDULE_MANAGER:
		active = "schedules"
	}

	count := func(query func() (int64, error)) int {
		if store == nil {
			return -1
		}
		n, err := query()
		if err != nil {
			if logger != nil {
				logger.Error("section nav count failed", "error", err)
			}
			return -1
		}
		return int(n)
	}

	tab := func(label, key, url string, count int) hb.TagInterface {
		class := "nav-link"
		if key == active {
			class += " active"
		}
		link := hb.Hyperlink().HTML(label).Href(url).Class(class)
		if count >= 0 {
			link.Child(hb.Span().
				Class("badge bg-secondary ms-2").
				HTML(cast.ToString(count)))
		}
		return hb.LI().Class("nav-item").Child(link)
	}

	ctx := r.Context()
	queueCount := count(func() (int64, error) {
		return store.TaskQueueCount(ctx, taskstore.TaskQueueQuery())
	})
	definitionCount := count(func() (int64, error) {
		return store.TaskDefinitionCount(ctx, taskstore.TaskDefinitionQuery())
	})
	scheduleCount := count(func() (int64, error) {
		return store.ScheduleCount(ctx, taskstore.NewScheduleQuery().SetLimit(0))
	})

	ulNav := hb.UL().
		Class("nav nav-pills p-1 justify-content-center flex-grow-1").
		Children([]hb.TagInterface{
			tab("Dashboard", "dashboard", linksHelper.Home(nil), -1),
			tab("Task Queue", "queue", linksHelper.TaskQueue(nil), queueCount),
			tab("Task Definitions", "definitions", linksHelper.TaskDefinitions(nil), definitionCount),
			tab("Schedules", "schedules", linksHelper.Schedules(nil), scheduleCount),
		})

	label := hb.Span().
		Class("fw-bolder text-uppercase ms-3 me-3").
		Text("Tasks").
		Child(hb.I().Class("bi bi-chevron-right ms-2"))

	return hb.Div().
		Class("card card-default mt-3 mb-4 bg-primary-subtle").
		Child(hb.Div().
			Class("card-body d-flex align-items-center").
			Style("padding: 2px;").
			Child(label).
			Child(ulNav))
}
