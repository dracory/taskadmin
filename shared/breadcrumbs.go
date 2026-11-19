package shared

import (
	"net/http"

	"github.com/dracory/hb"
	"github.com/samber/lo"
)

// Breadcrumbs renders the full breadcrumb trail for a page: the admin
// home link, the tasks admin link, then any page-specific items.
func Breadcrumbs(r *http.Request, pageBreadcrumbs []Breadcrumb) hb.TagInterface {
	adminHomeBreadcrumb := lo.
		If(AdminHomeURL(r) != "", Breadcrumb{
			Name: "Home",
			URL:  AdminHomeURL(r),
		}).
		Else(Breadcrumb{})

	breadcrumbItems := []Breadcrumb{
		adminHomeBreadcrumb,
		{
			Name: "Tasks",
			URL:  URL(r, CONTROLLER_HOME, nil),
		},
	}

	breadcrumbItems = append(breadcrumbItems, pageBreadcrumbs...)

	return BreadcrumbsUI(breadcrumbItems)
}

// BreadcrumbsUI renders a breadcrumb navigation from the given items
func BreadcrumbsUI(breadcrumbs []Breadcrumb) hb.TagInterface {
	ol := hb.OL().
		Class("breadcrumb").
		Style("border:none;")

	for _, breadcrumb := range breadcrumbs {
		link := hb.Hyperlink().
			HTML(breadcrumb.Name).
			Href(breadcrumb.URL)

		li := hb.LI().
			Class("breadcrumb-item").
			Child(link)

		ol.AddChild(li)
	}

	return hb.Nav().
		Attr("aria-label", "breadcrumb").
		Child(ol)
}
