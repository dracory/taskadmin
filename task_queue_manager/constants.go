package task_queue_manager

import (
	"context"
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
	CondFieldTaskName = "task_name"
	CondFieldDateFrom = "date_from"
	CondFieldDateTo   = "date_to"
)

// applyConditions maps filter conditions to the task queue query.
// Multiple conditions on the same field: last one wins. Date conditions
// accumulate. Returns false when a task_name filter matched no task
// definitions — the caller should return an empty result.
func applyConditions(ctx context.Context, store taskstore.StoreInterface, query taskstore.TaskQueueQueryInterface, conds []condition) bool {
	nameFilter := ""

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
		case CondFieldTaskName:
			nameFilter = val
		case CondFieldDateFrom:
			query.SetCreatedAtGte(val + " 00:00:00")
		case CondFieldDateTo:
			query.SetCreatedAtLte(val + " 23:59:59")
		}
	}

	if nameFilter != "" {
		ids := matchingDefinitionIDs(ctx, store, nameFilter)
		if len(ids) == 0 {
			return false
		}
		query.SetTaskIDIn(ids)
	}

	return true
}

// matchingDefinitionIDs returns the IDs of task definitions whose title or
// alias contains the name filter (case-insensitive).
func matchingDefinitionIDs(ctx context.Context, store taskstore.StoreInterface, nameFilter string) []string {
	definitions, err := store.TaskDefinitionList(ctx, taskstore.TaskDefinitionQuery())
	if err != nil || len(definitions) == 0 {
		return nil
	}

	needle := strings.ToLower(nameFilter)
	ids := []string{}
	for _, d := range definitions {
		if strings.Contains(strings.ToLower(d.GetTitle()), needle) ||
			strings.Contains(strings.ToLower(d.GetAlias()), needle) {
			ids = append(ids, d.GetID())
		}
	}
	return ids
}
