package shared

import "strings"

// IsJSON is a naive implementation for superficial, rough and fast
// checking for JSON
func IsJSON(str string) bool {
	if strings.HasPrefix(str, "{") && strings.HasSuffix(str, "}") {
		return true
	}

	if strings.HasPrefix(str, "[") && strings.HasSuffix(str, "]") {
		return true
	}

	return false
}
