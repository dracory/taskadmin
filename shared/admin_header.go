package shared

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/dracory/hb"
	"github.com/dracory/taskstore"
	"github.com/spf13/cast"
)

// AdminHeader renders the nav tabs (Dashboard / Task Queue / Task
// Definitions) with record counts, and highlights the active tab based
// on the request's "controller" parameter.
func AdminHeader(store taskstore.StoreInterface, logger *slog.Logger, r *http.Request) hb.TagInterface {
	queueCount, err := store.TaskQueueCount(context.Background(), taskstore.TaskQueueQuery())

	if err != nil {
		logger.Error(err.Error())
		queueCount = -1
	}

	taskCount, err := store.TaskDefinitionCount(context.Background(), taskstore.TaskDefinitionQuery())

	if err != nil {
		logger.Error(err.Error())
		taskCount = -1
	}

	active := r.URL.Query().Get("controller")
	if active == "" {
		active = CONTROLLER_HOME
	}

	tab := func(label, controller, key string, count int64) hb.TagInterface {
		class := "nav-link"
		if key == active {
			class += " active"
		}
		link := hb.Hyperlink().
			HTML(label).
			Href(URL(r, controller, nil)).
			Class(class)
		if count >= 0 {
			link.Child(hb.Span().
				Class("badge bg-secondary ms-2").
				HTML(cast.ToString(count)))
		}
		return hb.LI().Class("nav-item").Child(link)
	}

	return hb.UL().Class("nav nav-tabs mb-3").Children([]hb.TagInterface{
		tab("Dashboard", CONTROLLER_HOME, CONTROLLER_HOME, -1),
		tab("Task Queue", CONTROLLER_TASK_QUEUE_MANAGER, CONTROLLER_TASK_QUEUE_MANAGER, queueCount),
		tab("Task Definitions", CONTROLLER_TASK_DEFINITION_MANAGER, CONTROLLER_TASK_DEFINITION_MANAGER, taskCount),
	})
}
