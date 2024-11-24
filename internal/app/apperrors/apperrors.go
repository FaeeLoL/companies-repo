// Package apperrors defines application-specific error types and utilities.
package apperrors

import (
	"errors"
	"net/http"
)

type AppError struct {
	Code    int    // HTTP-status
	Message string // err msg
	Err     error
}

func (e *AppError) Error() string {
	return e.Message
}

func NewBadRequestError(message string) *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: message,
	}
}

func NewNotFoundError(message string) *AppError {
	return &AppError{
		Code:    http.StatusNotFound,
		Message: message,
	}
}

func (e *AppError) WithCause(err error) *AppError {
	e.Err = err
	return e
}

func NewInternalServerError(message string) *AppError {
	return &AppError{
		Code:    http.StatusInternalServerError,
		Message: message,
	}
}

func MapToAppError(err error) *AppError {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}

	return NewInternalServerError("unexpected error occurred")
}
