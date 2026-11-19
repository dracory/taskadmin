package shared

import "testing"

func Test_contextKeys(t *testing.T) {
	if KeyEndpoint != "endpoint" {
		t.Error("KeyEndpoint should be \"endpoint\"")
	}
	if KeyAdminHomeURL == "" {
		t.Error("KeyAdminHomeURL should not be empty")
	}
	if KeyTaskAdminURL == "" {
		t.Error("KeyTaskAdminURL should not be empty")
	}
}

func Test_defaultFavicon(t *testing.T) {
	if defaultFavicon == "" {
		t.Error("defaultFavicon should not be empty")
	}
	if len(defaultFavicon) < 100 {
		t.Error("defaultFavicon should be a base64 string (longer than 100 chars)")
	}
}

func Test_controllerConstants(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{"CONTROLLER_HOME", CONTROLLER_HOME},
		{"CONTROLLER_TASK_QUEUE_MANAGER", CONTROLLER_TASK_QUEUE_MANAGER},
		{"CONTROLLER_TASK_DEFINITION_MANAGER", CONTROLLER_TASK_DEFINITION_MANAGER},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value == "" {
				t.Errorf("%s should not be empty", tt.name)
			}
		})
	}
}
