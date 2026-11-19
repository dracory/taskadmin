package shared

import (
	"context"
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

	queueCount := count(func() (int64, error) {
		return store.TaskQueueCount(context.Background(), taskstore.TaskQueueQuery())
	})
	definitionCount := count(func() (int64, error) {
		return store.TaskDefinitionCount(context.Background(), taskstore.TaskDefinitionQuery())
	})

	return hb.UL().Class("nav nav-pills mb-3").Children([]hb.TagInterface{
		tab("Dashboard", "dashboard", linksHelper.Home(nil), -1),
		tab("Task Queue", "queue", linksHelper.TaskQueue(nil), queueCount),
		tab("Task Definitions", "definitions", linksHelper.TaskDefinitions(nil), definitionCount),
	})
}
