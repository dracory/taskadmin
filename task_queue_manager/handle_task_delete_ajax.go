package task_queue_manager

import (
	"encoding/json"
	"net/http"

	"github.com/dracory/api"
)

func (u *ui) handleTaskDelete(w http.ResponseWriter, r *http.Request) string {
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
		_, _ = w.Write([]byte(api.Error("Task ID is required").ToString()))
		return ""
	}

	if err := store.TaskQueueSoftDeleteByID(ctx, reqBody.ID); err != nil {
		u.Logger().Error("Failed to delete queued task", "error", err)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Failed to delete task").ToString()))
		return ""
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(api.Success("Task deleted successfully").ToString()))
	return ""
}

func (u *ui) handleTaskDeleteSelected(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	store := u.Store()
	if store == nil {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Task store not available").ToString()))
		return ""
	}

	var reqBody struct {
		IDs []string `json:"ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Invalid request body").ToString()))
		return ""
	}

	if len(reqBody.IDs) == 0 {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("No tasks selected").ToString()))
		return ""
	}

	for _, id := range reqBody.IDs {
		if err := store.TaskQueueSoftDeleteByID(ctx, id); err != nil {
			u.Logger().Error("Failed to delete queued task", "id", id, "error", err)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(api.Error("Failed to delete task " + id).ToString()))
			return ""
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(api.Success("Tasks deleted successfully").ToString()))
	return ""
}
