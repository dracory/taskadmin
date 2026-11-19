package schedule_manager

import (
	"encoding/json"
	"net/http"

	"github.com/dracory/api"
	"github.com/dracory/taskstore"
)

// scheduleToMap serializes a schedule for the frontend.
func scheduleToMap(schedule taskstore.ScheduleInterface, definitionTitle string) map[string]any {
	rule := schedule.GetRecurrenceRule()
	recurrence := map[string]any{}
	if rule != nil {
		daysOfWeek := []string{}
		for _, d := range rule.GetDaysOfWeek() {
			daysOfWeek = append(daysOfWeek, string(d))
		}
		monthsOfYear := []string{}
		for _, m := range rule.GetMonthsOfYear() {
			monthsOfYear = append(monthsOfYear, string(m))
		}
		recurrence = map[string]any{
			"frequency":    string(rule.GetFrequency()),
			"startsAt":     rule.GetStartsAt(),
			"endsAt":       rule.GetEndsAt(),
			"interval":     rule.GetInterval(),
			"daysOfWeek":   daysOfWeek,
			"daysOfMonth":  rule.GetDaysOfMonth(),
			"monthsOfYear": monthsOfYear,
		}
	}

	return map[string]any{
		"id":                  schedule.GetID(),
		"name":                schedule.GetName(),
		"description":         schedule.GetDescription(),
		"status":              schedule.GetStatus(),
		"queue_name":          schedule.GetQueueName(),
		"task_definition_id":  schedule.GetTaskDefinitionID(),
		"task_definition":     definitionTitle,
		"task_parameters":     schedule.GetTaskParameters(),
		"start_at":            schedule.GetStartAt(),
		"end_at":              schedule.GetEndAt(),
		"execution_count":     schedule.GetExecutionCount(),
		"max_execution_count": schedule.GetMaxExecutionCount(),
		"last_run_at":         schedule.GetLastRunAt(),
		"next_run_at":         schedule.GetNextRunAt(),
		"recurrence_rule":     recurrence,
		"created_at":          schedule.GetCreatedAt(),
		"updated_at":          schedule.GetUpdatedAt(),
	}
}

// definitionTitles returns a map of task definition ID to title.
func (u *ui) definitionTitles(r *http.Request) map[string]string {
	titles := map[string]string{}
	definitions, err := u.Store().TaskDefinitionList(r.Context(), taskstore.TaskDefinitionQuery().SetLimit(1000))
	if err != nil {
		u.Logger().Error("Failed to load task definitions", "error", err)
		return titles
	}
	for _, d := range definitions {
		titles[d.GetID()] = d.GetTitle()
	}
	return titles
}

func (u *ui) handleLoadSchedules(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	store := u.Store()
	if store == nil {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Task store not available").ToString()))
		return ""
	}

	var reqBody struct {
		Page             int    `json:"page"`
		PerPage          int    `json:"per_page"`
		Status           string `json:"status"`
		Name             string `json:"name"`
		QueueName        string `json:"queue_name"`
		TaskDefinitionID string `json:"task_definition_id"`
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

	query := taskstore.NewScheduleQuery().
		SetOffset(reqBody.Page * reqBody.PerPage).
		SetLimit(reqBody.PerPage)

	if reqBody.Status != "" {
		query.SetStatus(reqBody.Status)
	}
	if reqBody.Name != "" {
		query.SetName(reqBody.Name)
	}
	if reqBody.QueueName != "" {
		query.SetQueueName(reqBody.QueueName)
	}
	if reqBody.TaskDefinitionID != "" {
		query.SetTaskDefinitionID(reqBody.TaskDefinitionID)
	}

	schedules, err := store.ScheduleList(ctx, query)
	if err != nil {
		u.Logger().Error("Failed to load schedules", "error", err)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Failed to load schedules").ToString()))
		return ""
	}

	total, err := store.ScheduleCount(ctx, query)
	if err != nil {
		u.Logger().Error("Failed to count schedules", "error", err)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Failed to count schedules").ToString()))
		return ""
	}

	titles := u.definitionTitles(r)

	scheduleList := []map[string]any{}
	for _, schedule := range schedules {
		scheduleList = append(scheduleList, scheduleToMap(schedule, titles[schedule.GetTaskDefinitionID()]))
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(api.SuccessWithData("Schedules loaded successfully", map[string]any{
		"schedules": scheduleList,
		"total":     total,
		"page":      reqBody.Page,
		"per_page":  reqBody.PerPage,
	}).ToString()))
	return ""
}

// handleLoadDefinitions returns task definitions for the schedule form dropdown.
func (u *ui) handleLoadDefinitions(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	store := u.Store()
	if store == nil {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Task store not available").ToString()))
		return ""
	}

	definitions, err := store.TaskDefinitionList(ctx, taskstore.TaskDefinitionQuery().SetLimit(1000))
	if err != nil {
		u.Logger().Error("Failed to load task definitions", "error", err)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(api.Error("Failed to load task definitions").ToString()))
		return ""
	}

	definitionList := []map[string]any{}
	for _, definition := range definitions {
		definitionList = append(definitionList, map[string]any{
			"id":     definition.GetID(),
			"title":  definition.GetTitle(),
			"alias":  definition.GetAlias(),
			"status": definition.GetStatus(),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(api.SuccessWithData("Task definitions loaded successfully", map[string]any{
		"definitions": definitionList,
	}).ToString()))
	return ""
}
