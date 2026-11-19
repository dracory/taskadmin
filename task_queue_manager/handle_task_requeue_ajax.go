package task_queue_manager

import (
	"encoding/json"
	"net/http"

	"github.com/dracory/api"
	"github.com/dracory/taskstore"
	"github.com/spf13/cast"
)

func (u *ui) handleTaskRequeue(w http.ResponseWriter, r *http.Request) string {
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

	task, err := store.TaskDefinitionFindByID(ctx, queue.GetTaskID())
	if err != nil {
		u.Logger().Error("Failed to find task definition", "error", err)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Failed to find task").ToString()))
		return ""
	}
	if task == nil {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Task not found").ToString()))
		return ""
	}

	parametersAny := map[string]any{}
	if err := json.Unmarshal([]byte(reqBody.Parameters), &parametersAny); err != nil {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Task parameters is not valid JSON").ToString()))
		return ""
	}

	_, err = store.TaskDefinitionEnqueueByAlias(ctx, taskstore.DefaultQueueName, task.GetAlias(), cast.ToStringMap(parametersAny))
	if err != nil {
		u.Logger().Error("Failed to requeue task", "error", err)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Failed to requeue task").ToString()))
		return ""
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(api.Success("Task requeued successfully").ToString()))
	return ""
}
