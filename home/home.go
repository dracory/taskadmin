package home

import (
	"github.com/dracory/taskadmin/shared"
	"net/http"

	"github.com/dracory/cdn"
	"github.com/dracory/hb"
)

// UiInterface defines the home controller's UI interface
type UiInterface interface {
	shared.UiInterface
	Home(w http.ResponseWriter, r *http.Request)
}

// ui implements UiInterface
type ui struct {
	shared.UiBase
}

// UI creates a new home controller UI from the given config
func UI(config shared.UiConfig) UiInterface {
	return &ui{UiBase: shared.NewUiBase(config)}
}

// Home handles the home controller requests
func (u *ui) Home(w http.ResponseWriter, r *http.Request) {
	shared.WriteTag(w, u.ToTag(w, r))
}

func (c *ui) ToTag(w http.ResponseWriter, r *http.Request) hb.TagInterface {
	layout := c.Layout(r)
	data, errorMessage := c.prepareData(r)

	layout.SetTitle("Dashboard | Tasks")

	if errorMessage != "" {
		layout.SetBody(hb.Div().
			Class("alert alert-danger").
			Text(errorMessage).ToHTML())

		return hb.Raw(layout.Render(w, r))
	}

	htmxScript := `setTimeout(() => async function() {
		if (!window.htmx) {
			let script = document.createElement('script');
			document.head.appendChild(script);
			script.type = 'text/javascript';
			script.src = '` + cdn.Htmx_2_0_0() + `';
			await script.onload
		}
	}, 1000);`

	swalScript := `setTimeout(() => async function() {
		if (!window.Swal) {
			let script = document.createElement('script');
			document.head.appendChild(script);
			script.type = 'text/javascript';
			script.src = '` + cdn.Sweetalert2_11() + `';
			await script.onload
		}
	}, 1000);`

	layout.SetBody(c.page(data).ToHTML())
	layout.SetScripts([]string{htmxScript, swalScript})

	return hb.Raw(layout.Render(w, r))
}

func (controller *ui) page(data homeControllerData) hb.TagInterface {
	adminHeader := shared.AdminHeader(controller.Store(), controller.Logger(), data.request)
	breadcrumbs := shared.Breadcrumbs(data.request, []shared.Breadcrumb{
		{
			Name: "Dashboard",
			URL:  "",
		},
	})

	return hb.Div().
		Class("container py-4").
		Child(breadcrumbs).
		Child(shared.PageHeader("Tasks. Dashboard")).
		Child(adminHeader)
}

func (controller *ui) prepareData(r *http.Request) (data homeControllerData, errorMessage string) {
	data.request = r

	return data, ""
}

type homeControllerData struct {
	request *http.Request
}
