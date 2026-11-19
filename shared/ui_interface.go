package shared

import (
	"log/slog"
	"net/http"

	"github.com/dracory/taskstore"
)

// UiInterface defines the methods every subcontroller UI must implement.
type UiInterface interface {
	Store() taskstore.StoreInterface
	Logger() *slog.Logger

	// Layout returns the layout for the given request (FuncLayout
	// result, or DefaultWebpage when FuncLayout is nil).
	Layout(r *http.Request) Layout
}
