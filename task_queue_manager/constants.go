package task_queue_manager

import (
	"strings"

	"github.com/dracory/taskstore"
)

// condition is a single filter condition sent from the frontend.
type condition struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

// Condition field constants for the multi-filter system.
const (
	CondFieldStatus   = "status"
	CondFieldID       = "id"
	CondFieldTaskID   = "task_id"
	CondFieldDateFrom = "date_from"
	CondFieldDateTo   = "date_to"
)

// applyConditions maps filter conditions to the task queue query.
// Multiple conditions on the same field: last one wins. Date conditions
// accumulate.
func applyConditions(query taskstore.TaskQueueQueryInterface, conds []condition) {
	for _, c := range conds {
		val := strings.TrimSpace(c.Value)
		if val == "" {
			continue
		}
		switch c.Field {
		case CondFieldStatus:
			query.SetStatus(val)
		case CondFieldID:
			query.SetID(val)
		case CondFieldTaskID:
			query.SetTaskID(val)
		case CondFieldDateFrom:
			query.SetCreatedAtGte(val + " 00:00:00")
		case CondFieldDateTo:
			query.SetCreatedAtLte(val + " 23:59:59")
		}
	}
}
