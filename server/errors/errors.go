package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// ErrorCode represents error codes for categorization
type ErrorCode string

const (
	// Client errors
	ErrCodeInvalidInput   ErrorCode = "INVALID_INPUT"
	ErrCodeUnauthorized   ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden      ErrorCode = "FORBIDDEN"
	ErrCodeNotFound       ErrorCode = "NOT_FOUND"
	ErrCodeAlreadyExists  ErrorCode = "ALREADY_EXISTS"
	ErrCodeConflict       ErrorCode = "CONFLICT"
	ErrCodeTooManyRequest ErrorCode = "TOO_MANY_REQUESTS"

	// Server errors
	ErrCodeInternalServer  ErrorCode = "INTERNAL_SERVER_ERROR"
	ErrCodeDatabaseError   ErrorCode = "DATABASE_ERROR"
	ErrCodeExternalService ErrorCode = "EXTERNAL_SERVICE_ERROR"
	ErrCodeTokenGeneration ErrorCode = "TOKEN_GENERATION_ERROR"
)

// AppError represents application-specific error with structured information
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

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Err
}

// WithContext adds context information to the error
func (e *AppError) WithContext(key string, value any) *AppError {
	if e.Context == nil {
		e.Context = make(map[string]any)
	}
	e.Context[key] = value
	return e
}

// IsType checks if error is of specific type
func (e *AppError) IsType(code ErrorCode) bool {
	return e.Code == code
}

// Client error constructors
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

func NewAlreadyExists(message string) *AppError {
	return &AppError{
		Code:       ErrCodeAlreadyExists,
		Message:    message,
		HTTPStatus: http.StatusConflict,
	}
}

func NewTooManyRequests(message string) *AppError {
	return &AppError{
		Code:       ErrCodeTooManyRequest,
		Message:    message,
		HTTPStatus: http.StatusTooManyRequests,
	}
}

// Server error constructors
func NewInternalServerError(message string, err error) *AppError {
	return &AppError{
		Code:       ErrCodeInternalServer,
		Message:    message,
		HTTPStatus: http.StatusInternalServerError,
		Err:        err,
	}
}

func NewDatabaseError(message string, err error) *AppError {
	return &AppError{
		Code:       ErrCodeDatabaseError,
		Message:    message,
		HTTPStatus: http.StatusInternalServerError,
		Err:        err,
	}
}

func NewExternalServiceError(message string, err error) *AppError {
	return &AppError{
		Code:       ErrCodeExternalService,
		Message:    message,
		HTTPStatus: http.StatusInternalServerError,
		Err:        err,
	}
}

// Utility functions
func IsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}

func IsClientError(err error) bool {
	if appErr, ok := IsAppError(err); ok {
		return appErr.HTTPStatus >= 400 && appErr.HTTPStatus < 500
	}
	return false
}

func IsServerError(err error) bool {
	if appErr, ok := IsAppError(err); ok {
		return appErr.HTTPStatus >= 500
	}
	return false
}
