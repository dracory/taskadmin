package schedule_manager

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
	//go:embed schedules.html
	schedulesHTML string

	//go:embed schedules.js
	schedulesJS string
)

func (u *ui) renderPage(w http.ResponseWriter, r *http.Request) string {
	if u.Store() == nil {
		return shared.ErrorAlert("Task store is not initialized")
	}

	breadcrumbs := shared.Breadcrumbs(r, []shared.Breadcrumb{
		{Name: "Schedules", URL: ""},
	})

	linksHelper := shared.NewLinksFromRequest(r)
	urls := map[string]string{
		"loadSchedules":   linksHelper.Schedules(map[string]string{"action": actionLoadSchedules}),
		"loadDefinitions": linksHelper.Schedules(map[string]string{"action": actionLoadDefinitions}),
		"scheduleCreate":  linksHelper.Schedules(map[string]string{"action": actionScheduleCreate}),
		"scheduleUpdate":  linksHelper.Schedules(map[string]string{"action": actionScheduleUpdate}),
		"scheduleDelete":  linksHelper.Schedules(map[string]string{"action": actionScheduleDelete}),
		"scheduleDetails": linksHelper.Schedules(map[string]string{"action": actionScheduleDetails}),
	}

	urlsJSON, _ := json.Marshal(urls)
	urlsScript := hb.Script(fmt.Sprintf("window.scheduleUrls = %s;", string(urlsJSON)))

	content := hb.Div().
		Class("container").
		Class("py-4").
		Child(shared.SectionNav(r, u.Store(), u.Logger())).
		Child(breadcrumbs).
		Child(shared.PageHeader("Tasks. Schedules")).
		Child(urlsScript).
		Child(hb.Raw(schedulesHTML)).
		Child(hb.Script(schedulesJS))

	return u.Layout()(w, r, "Schedules | Tasks", content.ToHTML(), shared.LayoutOptions{
		ScriptURLs: []string{
			cdn.VueJs_3(),
			cdn.Notiflix_3_2_8(),
		},
		StyleURLs: []string{
			cdn.Notiflix_3_2_8_CSS(),
		},
	})
}
