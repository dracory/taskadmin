package shared

import (
	"log/slog"
	"net/http"

	"github.com/dracory/taskstore"
)

// UiBase is a base struct that implements shared.UiInterface.
// Subcontroller ui structs embed this to get the Store(), Logger(), and
// Layout() methods for free, following the shopadmin pattern.
type UiBase struct {
	StoreField      taskstore.StoreInterface
	LoggerField     *slog.Logger
	FuncLayoutField func(r *http.Request) Layout
}

// Store returns the configured task store
func (u UiBase) Store() taskstore.StoreInterface { return u.StoreField }

// Logger returns the configured logger
func (u UiBase) Logger() *slog.Logger { return u.LoggerField }

// Layout resolves the layout for the request
func (u UiBase) Layout(r *http.Request) Layout {
	if u.FuncLayoutField != nil {
		if layout := u.FuncLayoutField(r); layout != nil {
			return layout
		}
	}
	return DefaultWebpage("", "")
}

// NewUiBase creates a UiBase from a UiConfig
func NewUiBase(config UiConfig) UiBase {
	return UiBase{
		StoreField:      config.Store,
		LoggerField:     config.Logger,
		FuncLayoutField: config.FuncLayout,
	}
}
