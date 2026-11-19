package task_definition_manager

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/taskadmin/shared"
)

var (
	//go:embed definitions.html
	definitionsHTML string

	//go:embed definitions.js
	definitionsJS string
)

func (u *ui) renderPage(w http.ResponseWriter, r *http.Request) string {
	if u.Store() == nil {
		return shared.ErrorAlert("Task store is not initialized")
	}

	breadcrumbs := shared.Breadcrumbs(r, []shared.Breadcrumb{
		{Name: "Task Definitions", URL: ""},
	})

	linksHelper := shared.NewLinksFromRequest(r)
	urls := map[string]string{
		"loadDefinitions":   linksHelper.TaskDefinitions(map[string]string{"action": actionLoadDefinitions}),
		"definitionCreate":  linksHelper.TaskDefinitions(map[string]string{"action": actionDefinitionCreate}),
		"definitionUpdate":  linksHelper.TaskDefinitions(map[string]string{"action": actionDefinitionUpdate}),
		"definitionDelete":  linksHelper.TaskDefinitions(map[string]string{"action": actionDefinitionDelete}),
		"definitionDetails": linksHelper.TaskDefinitions(map[string]string{"action": actionDefinitionDetails}),
	}

	urlsJSON, _ := json.Marshal(urls)
	urlsScript := hb.Script(fmt.Sprintf("window.taskDefinitionUrls = %s;", string(urlsJSON)))

	content := hb.Div().
		Class("container").
		Class("py-4").
		Child(shared.SectionNav(r, u.Store(), u.Logger())).
		Child(breadcrumbs).
		Child(shared.PageHeader("Tasks. Definitions")).
		Child(urlsScript).
		Child(hb.Raw(definitionsHTML)).
		Child(hb.Script(definitionsJS))

	return u.Layout()(w, r, "Task Definitions | Tasks", content.ToHTML(), shared.LayoutOptions{
		ScriptURLs: []string{
			cdn.VueJs_3(),
			cdn.Notiflix_3_2_8(),
		},
		StyleURLs: []string{
			cdn.Notiflix_3_2_8_CSS(),
		},
	})
}
