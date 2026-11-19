package task_queue_manager

import (
	"encoding/json"
	"net/http"

	"github.com/dracory/api"
)

func (u *ui) handleTaskUpdateParameters(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	store := u.Store()
	if store == nil {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Task store not available").ToString()))
		return ""
	}

	var reqBody struct {
		ID         string `json:"id"`
		Parameters string `json:"parameters"`
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

	if reqBody.Parameters == "" {
		reqBody.Parameters = "{}"
	}

	parametersAny := map[string]any{}
	if err := json.Unmarshal([]byte(reqBody.Parameters), &parametersAny); err != nil {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Task parameters is not valid JSON").ToString()))
		return ""
	}

	queue, err := store.TaskQueueFindByID(ctx, reqBody.ID)
	if err != nil {
		u.Logger().Error("Failed to find queued task", "error", err)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Failed to find queued task").ToString()))
		return ""
	}
	if queue == nil {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Queued task not found").ToString()))
		return ""
	}

	queue.SetParameters(reqBody.Parameters)

	if err := store.TaskQueueUpdate(ctx, queue); err != nil {
		u.Logger().Error("Failed to update task parameters", "error", err)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Failed to update parameters").ToString()))
		return ""
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(api.Success("Parameters updated successfully").ToString()))
	return ""
}
