package task_definition_manager

import (
	"encoding/json"
	"net/http"

	"github.com/dracory/api"
	"github.com/dracory/taskstore"
)

func (u *ui) handleDefinitionCreate(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	store := u.Store()
	if store == nil {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Task store not available").ToString()))
		return ""
	}

	var reqBody struct {
		Title       string `json:"title"`
		Alias       string `json:"alias"`
		Status      string `json:"status"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Invalid request body").ToString()))
		return ""
	}

	if reqBody.Title == "" {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Title is required").ToString()))
		return ""
	}
	if reqBody.Alias == "" {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Alias is required").ToString()))
		return ""
	}
	if reqBody.Status == "" {
		reqBody.Status = taskstore.TaskDefinitionStatusActive
	}

	task := taskstore.NewTaskDefinition().
		SetTitle(reqBody.Title).
		SetAlias(reqBody.Alias).
		SetStatus(reqBody.Status).
		SetDescription(reqBody.Description)

	if err := store.TaskDefinitionCreate(ctx, task); err != nil {
		u.Logger().Error("Failed to create task definition", "error", err)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Failed to create task definition").ToString()))
		return ""
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(api.Success("Task definition created successfully").ToString()))
	return ""
}

func (u *ui) handleDefinitionUpdate(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	store := u.Store()
	if store == nil {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Task store not available").ToString()))
		return ""
	}

	var reqBody struct {
		ID          string `json:"id"`
		Title       string `json:"title"`
		Alias       string `json:"alias"`
		Status      string `json:"status"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Invalid request body").ToString()))
		return ""
	}

	if reqBody.ID == "" {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Task definition ID is required").ToString()))
		return ""
	}

	task, err := store.TaskDefinitionFindByID(ctx, reqBody.ID)
	if err != nil {
		u.Logger().Error("Failed to find task definition", "error", err)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Failed to find task definition").ToString()))
		return ""
	}
	if task == nil {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Task definition not found").ToString()))
		return ""
	}

	task.SetTitle(reqBody.Title).
		SetAlias(reqBody.Alias).
		SetStatus(reqBody.Status).
		SetDescription(reqBody.Description)

	if err := store.TaskDefinitionUpdate(ctx, task); err != nil {
		u.Logger().Error("Failed to update task definition", "error", err)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Failed to update task definition").ToString()))
		return ""
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(api.Success("Task definition updated successfully").ToString()))
	return ""
}

func (u *ui) handleDefinitionDelete(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	store := u.Store()
	if store == nil {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Task store not available").ToString()))
		return ""
	}

	var reqBody struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Invalid request body").ToString()))
		return ""
	}

	if reqBody.ID == "" {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Task definition ID is required").ToString()))
		return ""
	}

	if err := store.TaskDefinitionSoftDeleteByID(ctx, reqBody.ID); err != nil {
		u.Logger().Error("Failed to delete task definition", "error", err)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Failed to delete task definition").ToString()))
		return ""
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(api.Success("Task definition deleted successfully").ToString()))
	return ""
}

func (u *ui) handleDefinitionDetails(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	store := u.Store()
	if store == nil {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Task store not available").ToString()))
		return ""
	}

	var reqBody struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Invalid request body").ToString()))
		return ""
	}

	task, err := store.TaskDefinitionFindByID(ctx, reqBody.ID)
	if err != nil {
		u.Logger().Error("Failed to find task definition", "error", err)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Failed to find task definition").ToString()))
		return ""
	}
	if task == nil {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Task definition not found").ToString()))
		return ""
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(api.SuccessWithData("Task definition loaded successfully", map[string]any{
		"id":          task.GetID(),
		"title":       task.GetTitle(),
		"alias":       task.GetAlias(),
		"description": task.GetDescription(),
		"status":      task.GetStatus(),
		"created_at":  task.GetCreatedAt(),
		"updated_at":  task.GetUpdatedAt(),
	}).ToString()))
	return ""
}
