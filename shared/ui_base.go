package shared

import (
	"log/slog"

	"github.com/dracory/taskstore"
)

// UiBase is a base struct that implements shared.UiInterface.
// Subcontroller ui structs embed this to get the Store(), Logger(), and
// Layout() methods for free, following the shopadmin pattern.
type UiBase struct {
	StoreField  taskstore.StoreInterface
	LoggerField *slog.Logger
	LayoutField FuncLayout
}

// Store returns the configured task store
func (u UiBase) Store() taskstore.StoreInterface { return u.StoreField }

// Logger returns the configured logger
func (u UiBase) Logger() *slog.Logger { return u.LoggerField }

// Layout returns the configured layout function
func (u UiBase) Layout() FuncLayout { return u.LayoutField }

// NewUiBase creates a UiBase from a UiConfig
func NewUiBase(config UiConfig) UiBase {
	return UiBase{
		StoreField:  config.Store,
		LoggerField: config.Logger,
		LayoutField: config.Layout,
	}
}
