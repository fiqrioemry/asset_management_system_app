// ==================== errors/errors.go ====================
package errors

import (
	"errors"
	"fmt"
	"net/http"
)

type ErrorCode string

const (
	// Client errors (4xx)
	ErrCodeInvalidInput    ErrorCode = "INVALID_INPUT"     // 400
	ErrCodeUnauthorized    ErrorCode = "UNAUTHORIZED"      // 401
	ErrCodeForbidden       ErrorCode = "FORBIDDEN"         // 403
	ErrCodeNotFound        ErrorCode = "NOT_FOUND"         // 404
	ErrCodeConflict        ErrorCode = "CONFLICT"          // 409
	ErrCodeRequestTooLarge ErrorCode = "REQUEST_TOO_LARGE" // 413
	ErrCodeTooManyRequest  ErrorCode = "TOO_MANY_REQUESTS" // 429

	// Server errors (5xx)
	ErrCodeInternalServer  ErrorCode = "INTERNAL_SERVER_ERROR"
	ErrCodeDatabaseError   ErrorCode = "DATABASE_ERROR"
	ErrCodeExternalService ErrorCode = "EXTERNAL_SERVICE_ERROR"
)

type AppError struct {
	Code       ErrorCode      `json:"code"`
	Message    string         `json:"message"`
	HTTPStatus int            `json:"-"`
	Err        error          `json:"-"`
	Context    map[string]any `json:"context,omitempty"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) WithContext(key string, value any) *AppError {
	if e.Context == nil {
		e.Context = make(map[string]any)
	}
	e.Context[key] = value
	return e
}

func NewTooManyRequests(message string) *AppError {
	return &AppError{
		Code:       ErrCodeTooManyRequest,
		Message:    message,
		HTTPStatus: http.StatusTooManyRequests,
	}
}

func NewBadRequest(message string) *AppError {
	return &AppError{
		Code:       ErrCodeInvalidInput,
		Message:    message,
		HTTPStatus: http.StatusBadRequest,
	}
}

func NewUnauthorized(message string) *AppError {
	return &AppError{
		Code:       ErrCodeUnauthorized,
		Message:    message,
		HTTPStatus: http.StatusUnauthorized,
	}
}

func NewForbidden(message string) *AppError {
	return &AppError{
		Code:       ErrCodeForbidden,
		Message:    message,
		HTTPStatus: http.StatusForbidden,
	}
}

func NewNotFound(message string) *AppError {
	return &AppError{
		Code:       ErrCodeNotFound,
		Message:    message,
		HTTPStatus: http.StatusNotFound,
	}
}

func NewConflict(message string) *AppError {
	return &AppError{
		Code:       ErrCodeConflict,
		Message:    message,
		HTTPStatus: http.StatusConflict,
	}
}

func NewInternalServerError(message string, err error) *AppError {
	return &AppError{
		Code:       ErrCodeInternalServer,
		Message:    message,
		HTTPStatus: http.StatusInternalServerError,
		Err:        err,
	}
}

func IsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}

func NewRequestTooLarge(message string) *AppError {
	return &AppError{
		Code:       ErrCodeRequestTooLarge,
		Message:    message,
		HTTPStatus: http.StatusRequestEntityTooLarge,
	}
}

func IsServerError(err error) bool {
	if appErr, ok := IsAppError(err); ok {
		return appErr.HTTPStatus >= 500
	}
	return false
}
