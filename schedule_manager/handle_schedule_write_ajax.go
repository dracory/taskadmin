package schedule_manager

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/dracory/api"
	"github.com/dracory/taskstore"
	"github.com/dromara/carbon/v2"
)

// scheduleRequest is the JSON payload sent by the schedule form.
type scheduleRequest struct {
	ID                string         `json:"id"`
	Name              string         `json:"name"`
	Description       string         `json:"description"`
	Status            string         `json:"status"`
	QueueName         string         `json:"queue_name"`
	TaskDefinitionID  string         `json:"task_definition_id"`
	TaskParameters    map[string]any `json:"task_parameters"`
	StartAt           string         `json:"start_at"`
	EndAt             string         `json:"end_at"`
	MaxExecutionCount int            `json:"max_execution_count"`

	// Recurrence rule fields
	Frequency    string   `json:"frequency"`
	Interval     int      `json:"interval"`
	DaysOfWeek   []string `json:"days_of_week"`
	DaysOfMonth  []int    `json:"days_of_month"`
	MonthsOfYear []string `json:"months_of_year"`
}

// normalizeDateTime converts an HTML datetime-local value ("2006-01-02T15:04")
// to the store's "2006-01-02 15:04:05" format. Empty input returns "".
func normalizeDateTime(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	value = strings.ReplaceAll(value, "T", " ")
	if len(value) == 16 {
		value += ":00"
	}
	return value
}

// jsonError writes a JSON error response and returns "".
func jsonError(w http.ResponseWriter, message string) string {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(api.Error(message).ToString()))
	return ""
}

// jsonSuccess writes a JSON success response and returns "".
func jsonSuccess(w http.ResponseWriter, message string) string {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(api.Success(message).ToString()))
	return ""
}

// applyRequest validates the request and applies it to the schedule,
// rebuilding the recurrence rule and recalculating next_run_at.
func applyRequest(schedule taskstore.ScheduleInterface, req scheduleRequest) error {
	status := req.Status
	if status == "" {
		status = ScheduleStatusDraft
	}

	startAt := normalizeDateTime(req.StartAt)
	if startAt == "" {
		startAt = carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)
	}
	endAt := normalizeDateTime(req.EndAt)
	if endAt == "" {
		endAt = taskstore.MAX_DATETIME
	}

	interval := req.Interval
	if interval <= 0 {
		interval = 1
	}

	daysOfWeek := make([]taskstore.DayOfWeek, 0, len(req.DaysOfWeek))
	for _, d := range req.DaysOfWeek {
		daysOfWeek = append(daysOfWeek, taskstore.DayOfWeek(strings.ToLower(strings.TrimSpace(d))))
	}
	monthsOfYear := make([]taskstore.MonthOfYear, 0, len(req.MonthsOfYear))
	for _, m := range req.MonthsOfYear {
		monthsOfYear = append(monthsOfYear, taskstore.MonthOfYear(strings.ToUpper(strings.TrimSpace(m))))
	}

	rule := taskstore.NewRecurrenceRule().
		SetFrequency(taskstore.Frequency(req.Frequency)).
		SetInterval(interval).
		SetStartsAt(startAt).
		SetEndsAt(endAt).
		SetDaysOfWeek(daysOfWeek).
		SetDaysOfMonth(req.DaysOfMonth).
		SetMonthsOfYear(monthsOfYear)

	if req.TaskParameters == nil {
		req.TaskParameters = map[string]any{}
	}

	schedule.
		SetName(req.Name).
		SetDescription(req.Description).
		SetStatus(status).
		SetQueueName(req.QueueName).
		SetTaskDefinitionID(req.TaskDefinitionID).
		SetTaskParameters(req.TaskParameters).
		SetStartAt(startAt).
		SetEndAt(endAt).
		SetMaxExecutionCount(req.MaxExecutionCount).
		SetRecurrenceRule(rule).
		UpdateNextRunAt()

	return nil
}

func (u *ui) handleScheduleCreate(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	store := u.Store()
	if store == nil {
		return jsonError(w, "Task store not available")
	}

	var reqBody scheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		return jsonError(w, "Invalid request body")
	}

	if strings.TrimSpace(reqBody.Name) == "" {
		return jsonError(w, "Name is required")
	}
	if strings.TrimSpace(reqBody.TaskDefinitionID) == "" {
		return jsonError(w, "Task definition is required")
	}

	schedule := taskstore.NewSchedule()
	if err := applyRequest(schedule, reqBody); err != nil {
		return jsonError(w, err.Error())
	}

	if err := store.ScheduleCreate(ctx, schedule); err != nil {
		u.Logger().Error("Failed to create schedule", "error", err)
		return jsonError(w, "Failed to create schedule")
	}

	return jsonSuccess(w, "Schedule created successfully")
}

func (u *ui) handleScheduleUpdate(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	store := u.Store()
	if store == nil {
		return jsonError(w, "Task store not available")
	}

	var reqBody scheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		return jsonError(w, "Invalid request body")
	}

	if reqBody.ID == "" {
		return jsonError(w, "Schedule ID is required")
	}
	if strings.TrimSpace(reqBody.Name) == "" {
		return jsonError(w, "Name is required")
	}

	schedule, err := store.ScheduleFindByID(ctx, reqBody.ID)
	if err != nil {
		u.Logger().Error("Failed to find schedule", "error", err)
		return jsonError(w, "Failed to find schedule")
	}
	if schedule == nil {
		return jsonError(w, "Schedule not found")
	}

	if err := applyRequest(schedule, reqBody); err != nil {
		return jsonError(w, err.Error())
	}

	if err := store.ScheduleUpdate(ctx, schedule); err != nil {
		u.Logger().Error("Failed to update schedule", "error", err)
		return jsonError(w, "Failed to update schedule")
	}

	return jsonSuccess(w, "Schedule updated successfully")
}

func (u *ui) handleScheduleDelete(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	store := u.Store()
	if store == nil {
		return jsonError(w, "Task store not available")
	}

	var reqBody struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		return jsonError(w, "Invalid request body")
	}

	if reqBody.ID == "" {
		return jsonError(w, "Schedule ID is required")
	}

	if err := store.ScheduleSoftDeleteByID(ctx, reqBody.ID); err != nil {
		u.Logger().Error("Failed to delete schedule", "error", err)
		return jsonError(w, "Failed to delete schedule")
	}

	return jsonSuccess(w, "Schedule deleted successfully")
}

func (u *ui) handleScheduleDetails(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	store := u.Store()
	if store == nil {
		return jsonError(w, "Task store not available")
	}

	var reqBody struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		return jsonError(w, "Invalid request body")
	}

	schedule, err := store.ScheduleFindByID(ctx, reqBody.ID)
	if err != nil {
		u.Logger().Error("Failed to find schedule", "error", err)
		return jsonError(w, "Failed to find schedule")
	}
	if schedule == nil {
		return jsonError(w, "Schedule not found")
	}

	titles := u.definitionTitles(r)

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(api.SuccessWithData("Schedule loaded successfully",
		scheduleToMap(schedule, titles[schedule.GetTaskDefinitionID()])).ToString()))
	return ""
}
