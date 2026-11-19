package taskadmin

import "errors"

var (
	ErrStoreRequired  = errors.New("store is required")
	ErrLoggerRequired = errors.New("logger is required")
)
