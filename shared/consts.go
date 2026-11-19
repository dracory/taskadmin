package shared

// Context keys for config values injected by Handle()
const KeyEndpoint = "endpoint"
const KeyAdminHomeURL = "admin_home_url"
const KeyTaskAdminURL = "task_admin_url"

// Controller names used in the ?controller= query parameter
const (
	CONTROLLER_HOME                    = "home"
	CONTROLLER_TASK_QUEUE_MANAGER      = "task-queue-manager"
	CONTROLLER_TASK_DEFINITION_MANAGER = "task-definition-manager"
)
