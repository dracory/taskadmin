package shared

// Context keys for config values injected by Handle()
const KeyEndpoint = "endpoint"
const KeyAdminHomeURL = "admin_home_url"
const KeyTaskAdminURL = "task_admin_url"

// Controller names used in the ?controller= query parameter
const (
	CONTROLLER_HOME = "home"

	CONTROLLER_TASK_QUEUE_CREATE       = "task-queue-create"
	CONTROLLER_TASK_QUEUE_DELETE       = "task-queue-delete"
	CONTROLLER_TASK_QUEUE_DETAILS      = "task-queue-details"
	CONTROLLER_TASK_QUEUE_MANAGER      = "task-queue-manager"
	CONTROLLER_TASK_QUEUE_PARAMETERS   = "task-queue-parameters"
	CONTROLLER_TASK_QUEUE_REQUEUE      = "task-queue-requeue"
	CONTROLLER_TASK_QUEUE_TASK_RESTART = "task-queue-task-restart"

	CONTROLLER_TASK_DEFINITION_CREATE  = "task-definition-create"
	CONTROLLER_TASK_DEFINITION_MANAGER = "task-definition-manager"
	CONTROLLER_TASK_DEFINITION_UPDATE  = "task-definition-update"
	CONTROLLER_TASK_DEFINITION_DELETE  = "task-definition-delete"
)

// Form field names
const (
	FIELD_PARAMETERS  = "parameters"
	FIELD_QUEUE_ID    = "queue_id"
	FIELD_TASK_ID     = "task_id"
	FIELD_STATUS      = "status"
	FIELD_TITLE       = "title"
	FIELD_ALIAS       = "alias"
	FIELD_DESCRIPTION = "description"
	FIELD_DETAILS     = "details"

	FIELD_FILTER_QUEUE_ID     = "filter_queue_id"
	FIELD_FILTER_STATUS       = "filter_status"
	FIELD_FILTER_NAME         = "filter_name"
	FIELD_FILTER_CREATED_FROM = "filter_created_from"
	FIELD_FILTER_CREATED_TO   = "filter_created_to"
	FIELD_FILTER_TASK_ID      = "filter_task_id"
)

// Modal / form action names
const (
	ACTION_MODAL_QUEUED_TASK_FILTER_SHOW       = "modal-queued-task-filter-show"
	ACTION_MODAL_QUEUED_TASK_RESTART_SHOW      = "modal-queue-task-restart-show"
	ACTION_MODAL_QUEUED_TASK_RESTART_SUBMITTED = "modal-queue-task-restart-submitted"
	ACTION_MODAL_QUEUE_FILTER_SHOW             = "modal_queue_filter_show"
)
