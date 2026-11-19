package shared

import (
	"log/slog"

	"github.com/dracory/taskstore"
)

// UiInterface defines the methods every subcontroller UI must implement.
type UiInterface interface {
	Store() taskstore.StoreInterface
	Logger() *slog.Logger
	Layout() FuncLayout
}
