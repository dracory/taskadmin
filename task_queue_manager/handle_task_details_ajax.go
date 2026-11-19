package task_queue_manager

import (
	"encoding/json"
	"net/http"

	"github.com/dracory/api"
)

func (u *ui) handleTaskDetails(w http.ResponseWriter, r *http.Request) string {
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

	taskTitle := ""
	if task, err := store.TaskDefinitionFindByID(ctx, queue.GetTaskID()); err == nil && task != nil {
		taskTitle = task.GetTitle()
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(api.SuccessWithData("Task loaded successfully", map[string]any{
		"id":           queue.GetID(),
		"task_id":      queue.GetTaskID(),
		"task_title":   taskTitle,
		"queue_name":   queue.GetQueueName(),
		"status":       queue.GetStatus(),
		"parameters":   queue.GetParameters(),
		"output":       queue.GetOutput(),
		"details":      queue.GetDetails(),
		"attempts":     queue.GetAttempts(),
		"started_at":   queue.GetStartedAt(),
		"completed_at": queue.GetCompletedAt(),
		"created_at":   queue.GetCreatedAt(),
		"updated_at":   queue.GetUpdatedAt(),
	}).ToString()))
	return ""
}
