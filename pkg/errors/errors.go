package errors

import "errors"

// Sentinel errors untuk API
var (
	ErrInvalidRequest     = errors.New("invalid request")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrNotFound           = errors.New("not found")
	ErrConflict           = errors.New("conflict")
	ErrInternalServer     = errors.New("internal server error")
	ErrServiceUnavailable = errors.New("service unavailable")
	ErrBadGateway         = errors.New("bad gateway")
)

// APIError wraps error dengan HTTP status code
type APIError struct {
	StatusCode int
	Message    string
	Err        error
}

// Error implements error interface
func (e *APIError) Error() string {
	return e.Message
}

// NewAPIError creates new API error
func NewAPIError(statusCode int, message string, err error) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Message:    message,
		Err:        err,
	}
}
