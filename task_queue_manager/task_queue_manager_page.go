package task_queue_manager

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/taskadmin/shared"
	"github.com/dracory/taskstore"
)

var (
	//go:embed tasks.html
	tasksHTML string

	//go:embed tasks.js
	tasksJS string
)

func (u *ui) renderPage(w http.ResponseWriter, r *http.Request) string {
	if u.Store() == nil {
		return shared.ErrorAlert("Task store is not initialized")
	}

	breadcrumbs := shared.Breadcrumbs(r, []shared.Breadcrumb{
		{Name: "Task Queue", URL: ""},
	})

	linksHelper := shared.NewLinksFromRequest(r)
	urls := map[string]string{
		"loadTasks":          linksHelper.TaskQueue(map[string]string{"action": actionLoadTasks}),
		"taskEnqueue":        linksHelper.TaskQueue(map[string]string{"action": actionTaskEnqueue}),
		"taskDelete":         linksHelper.TaskQueue(map[string]string{"action": actionTaskDelete}),
		"taskDeleteSelected": linksHelper.TaskQueue(map[string]string{"action": actionTaskDeleteSelected}),
		"taskRequeue":        linksHelper.TaskQueue(map[string]string{"action": actionTaskRequeue}),
		"taskRestart":        linksHelper.TaskQueue(map[string]string{"action": actionTaskRestart}),
		"taskDetails":        linksHelper.TaskQueue(map[string]string{"action": actionTaskDetails}),
		"taskUpdateParams":   linksHelper.TaskQueue(map[string]string{"action": actionTaskUpdateParams}),
	}

	urlsJSON, _ := json.Marshal(urls)
	urlsScript := hb.Script(fmt.Sprintf("window.taskQueueUrls = %s;", string(urlsJSON)))

	definitions, err := u.Store().TaskDefinitionList(r.Context(), taskstore.TaskDefinitionQuery().
		SetStatus(taskstore.TaskDefinitionStatusActive).
		SetOrderBy(taskstore.COLUMN_TITLE).
		SetSortOrder(taskstore.ASC))
	definitionItems := []map[string]any{}
	if err != nil {
		u.Logger().Error("Failed to load task definitions", "error", err)
	} else {
		for _, definition := range definitions {
			definitionItems = append(definitionItems, map[string]any{
				"id":    definition.GetID(),
				"title": definition.GetTitle(),
				"alias": definition.GetAlias(),
			})
		}
	}
	definitionsJSON, _ := json.Marshal(definitionItems)
	definitionsScript := hb.Script(fmt.Sprintf("window.taskDefinitions = %s;", string(definitionsJSON)))

	content := hb.Div().
		Class("container").
		Class("py-4").
		Child(shared.SectionNav(r, u.Store(), u.Logger())).
		Child(breadcrumbs).
		Child(shared.PageHeader("Tasks. Queue")).
		Child(urlsScript).
		Child(definitionsScript).
		Child(hb.Raw(tasksHTML)).
		Child(hb.Script(tasksJS))

	return u.Layout()(w, r, "Queue Manager | Tasks", content.ToHTML(), shared.LayoutOptions{
		ScriptURLs: []string{
			cdn.VueJs_3(),
			cdn.Notiflix_3_2_8(),
		},
		StyleURLs: []string{
			cdn.Notiflix_3_2_8_CSS(),
		},
	})
}
