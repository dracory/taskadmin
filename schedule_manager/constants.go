package schedule_manager

// Schedule status values. taskstore does not export constants for these;
// valid values are "draft" (default), "active", "inactive".
const (
	ScheduleStatusDraft    = "draft"
	ScheduleStatusActive   = "active"
	ScheduleStatusInactive = "inactive"
)
