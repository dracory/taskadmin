package shared

import (
	"strings"

	"github.com/dracory/hb"
	"github.com/samber/lo"
)

// SortingIndicator returns the sorting indicator (up/down arrow) for a
// column. This is a standalone utility for reusability across controllers.
func SortingIndicator(columnName, sortByColumnName, sortOrder string) hb.TagInterface {
	isSelected := strings.EqualFold(sortByColumnName, columnName)

	direction := lo.If(isSelected && sortOrder == "asc", "up").
		ElseIf(isSelected && sortOrder == "desc", "down").
		Else("none")

	return hb.Span().
		Class("sorting").
		HTMLIf(direction == "up", "&#8595;").
		HTMLIf(direction == "down", "&#8593;")
}
