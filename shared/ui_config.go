package shared

import (
	"log/slog"
	"net/http"

	"github.com/dracory/taskstore"
)

// UiConfig holds the dependencies passed to subcontroller UI factories.
//
// FuncLayout returns a fresh Layout per request; when nil,
// DefaultWebpage is used instead.
type UiConfig struct {
	Store      taskstore.StoreInterface
	Logger     *slog.Logger
	FuncLayout func(r *http.Request) Layout
}
