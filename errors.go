package heimdall

import "fmt"

// Error codes for config operations.
const (
	ErrConfigNotFound   = -32002
	ErrConfigValidation = -32003
	ErrPermissionDenied = -32004
	ErrInvalidParams    = -32602
)

// ConfigError represents a configuration operation error.
type ConfigError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *ConfigError) Error() string {
	return fmt.Sprintf("heimdall error %d: %s", e.Code, e.Message)
}
