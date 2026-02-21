package domain

import "fmt"

// ErrorCode defines error types in the system.
type ErrorCode string

const (
	ErrSearchFailed   ErrorCode = "SEARCH_FAILED"
	ErrDownloadFailed ErrorCode = "DOWNLOAD_FAILED"
	ErrJobQueueFull   ErrorCode = "JOB_QUEUE_FULL"
	ErrInvalidInput   ErrorCode = "INVALID_INPUT"
	ErrNotFound       ErrorCode = "NOT_FOUND"
	ErrInternalServer ErrorCode = "INTERNAL_SERVER_ERROR"
	ErrUnauthorized   ErrorCode = "UNAUTHORIZED"
	ErrConflict       ErrorCode = "CONFLICT"
)

// AppError is the standard error type for the application.
type AppError struct {
	Code    ErrorCode              `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
	Err     error                  `json:"-"` // Internal error, not exposed to users
}

// Error implements the error interface.
func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// NewAppError creates a new AppError.
func NewAppError(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Details: make(map[string]interface{}),
	}
}

// WithError wraps an underlying error.
func (e *AppError) WithError(err error) *AppError {
	e.Err = err
	return e
}

// WithDetails adds contextual details.
func (e *AppError) WithDetails(details map[string]interface{}) *AppError {
	e.Details = details
	return e
}
