package shared

import (
	"log/slog"
	"net/http"

	"github.com/dracory/taskstore"
)

// FuncLayout is the layout function signature used by subcontrollers.
// It receives the full page (title + body HTML + optional assets) and
// returns the rendered page. When w is non-nil implementations may
// write directly to the response and return "".
type FuncLayout func(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options LayoutOptions) string

// UiConfig holds the dependencies passed to subcontroller UI factories.
type UiConfig struct {
	Store  taskstore.StoreInterface
	Logger *slog.Logger
	Layout FuncLayout
}
