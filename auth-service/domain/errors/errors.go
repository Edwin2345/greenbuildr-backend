package errors

import (
	"fmt"
	"net/http"
)

// DomainError represents a domain-level error with metadata
type DomainError struct {
	Code       string // Machine-readable error code
	Message    string // Human-readable message
	StatusCode int    // HTTP status code
	Err        error  // Underlying error
}

// Error implements the error interface
func (e *DomainError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Domain-specific errors
var (
	ErrUserNotFound = &DomainError{
		Code:       "USER_NOT_FOUND",
		Message:    "User not found",
		StatusCode: http.StatusNotFound,
	}

	ErrEmailAlreadyExists = &DomainError{
		Code:       "EMAIL_ALREADY_EXISTS",
		Message:    "Email already registered",
		StatusCode: http.StatusConflict,
	}

	ErrInvalidPassword = &DomainError{
		Code:       "INVALID_PASSWORD",
		Message:    "Invalid password",
		StatusCode: http.StatusUnauthorized,
	}

	ErrInvalidEmail = &DomainError{
		Code:       "INVALID_EMAIL",
		Message:    "Invalid email format",
		StatusCode: http.StatusBadRequest,
	}

	ErrPasswordTooShort = &DomainError{
		Code:       "PASSWORD_TOO_SHORT",
		Message:    "Password must be at least 8 characters",
		StatusCode: http.StatusBadRequest,
	}

	ErrInvalidToken = &DomainError{
		Code:       "INVALID_TOKEN",
		Message:    "Invalid or expired token",
		StatusCode: http.StatusUnauthorized,
	}

	ErrInternalError = &DomainError{
		Code:       "INTERNAL_ERROR",
		Message:    "Internal server error",
		StatusCode: http.StatusInternalServerError,
	}
)

// NewDomainError creates a new domain error
func NewDomainError(code, message string, statusCode int, err error) *DomainError {
	return &DomainError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
		Err:        err,
	}
}
