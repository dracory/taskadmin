package task_queue_manager

import (
	"encoding/json"
	"net/http"

	"github.com/dracory/api"
	"github.com/dracory/taskstore"
)

func (u *ui) handleLoadTasks(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	store := u.Store()
	if store == nil {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Task store not available").ToString()))
		return ""
	}

	var reqBody struct {
		Page        int         `json:"page"`
		PerPage     int         `json:"per_page"`
		SortBy      string      `json:"sort_by"`
		Sort        string      `json:"sort"`
		Status      string      `json:"status"`       // legacy single filter
		CreatedFrom string      `json:"created_from"` // legacy single filter
		CreatedTo   string      `json:"created_to"`   // legacy single filter
		Conditions  []condition `json:"conditions"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Invalid request body").ToString()))
		return ""
	}

	if reqBody.Page < 0 {
		reqBody.Page = 0
	}
	if reqBody.PerPage <= 0 {
		reqBody.PerPage = 10
	}
	if reqBody.SortBy == "" {
		reqBody.SortBy = taskstore.COLUMN_CREATED_AT
	}
	if reqBody.Sort == "" {
		reqBody.Sort = taskstore.DESC
	}

	query := taskstore.TaskQueueQuery().
		SetOffset(reqBody.Page * reqBody.PerPage).
		SetLimit(reqBody.PerPage).
		SetOrderBy(reqBody.SortBy).
		SetSortOrder(reqBody.Sort)

	if reqBody.Status != "" {
		query.SetStatus(reqBody.Status)
	}
	if reqBody.CreatedFrom != "" {
		query.SetCreatedAtGte(reqBody.CreatedFrom + " 00:00:00")
	}
	if reqBody.CreatedTo != "" {
		query.SetCreatedAtLte(reqBody.CreatedTo + " 23:59:59")
	}

	if !applyConditions(ctx, store, query, reqBody.Conditions) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.SuccessWithData("Tasks loaded successfully", map[string]any{
			"tasks":    []map[string]any{},
			"total":    0,
			"page":     reqBody.Page,
			"per_page": reqBody.PerPage,
		}).ToString()))
		return ""
	}

	tasks, err := store.TaskQueueList(ctx, query)
	if err != nil {
		u.Logger().Error("Failed to load queued tasks", "error", err)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Failed to load tasks").ToString()))
		return ""
	}

	total, err := store.TaskQueueCount(ctx, query)
	if err != nil {
		u.Logger().Error("Failed to count queued tasks", "error", err)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Failed to count tasks").ToString()))
		return ""
	}

	definitions := map[string]taskstore.TaskDefinitionInterface{}
	definitionFor := func(taskID string) taskstore.TaskDefinitionInterface {
		if taskID == "" {
			return nil
		}
		if definition, ok := definitions[taskID]; ok {
			return definition
		}
		definition, err := store.TaskDefinitionFindByID(ctx, taskID)
		if err != nil {
			u.Logger().Error("Failed to load task definition", "task_id", taskID, "error", err)
		}
		definitions[taskID] = definition
		return definition
	}

	taskList := []map[string]any{}
	for _, task := range tasks {
		taskTitle := ""
		taskAlias := ""
		if definition := definitionFor(task.GetTaskID()); definition != nil {
			taskTitle = definition.GetTitle()
			taskAlias = definition.GetAlias()
		}
		taskList = append(taskList, map[string]any{
			"id":         task.GetID(),
			"task_id":    task.GetTaskID(),
			"task_title": taskTitle,
			"task_alias": taskAlias,
			"status":     task.GetStatus(),
			"attempts":   task.GetAttempts(),
			"created_at": task.GetCreatedAt(),
			"updated_at": task.GetUpdatedAt(),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(api.SuccessWithData("Tasks loaded successfully", map[string]any{
		"tasks":    taskList,
		"total":    total,
		"page":     reqBody.Page,
		"per_page": reqBody.PerPage,
	}).ToString()))
	return ""
}
