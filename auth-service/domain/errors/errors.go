package errors

import (
	"net/http"

	"github.com/greenbuildr/shared/domainerrors"
)

type DomainError = domainerrors.DomainError

var (
	//ErrUserNotFound is internal use only
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

	ErrInvalidCredentials = &DomainError{
		Code:       "INVALID_CREDENTIALS",
		Message:    "Invalid email or password",
		StatusCode: http.StatusUnauthorized,
	}

	ErrEmailNotVerified = &DomainError{
		Code:       "EMAIL_NOT_VERIFIED",
		Message:    "Please verify your email before logging in",
		StatusCode: http.StatusForbidden,
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
