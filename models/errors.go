package models

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

var (
	ErrUnauthorized      = NewAppError(http.StatusUnauthorized, "Unauthorized", nil)
	ErrInvalidFileType   = NewAppError(http.StatusBadRequest, "Invalid file type", nil)
	ErrFileTooLarge      = NewAppError(http.StatusBadRequest, "File too large", nil)
	ErrFileNotFound      = NewAppError(http.StatusNotFound, "File not found", nil)
	ErrInternalServer    = NewAppError(http.StatusInternalServerError, "Internal server error", nil)
	ErrBadRequest        = NewAppError(http.StatusBadRequest, "Bad request", nil)
	ErrRateLimitExceeded = NewAppError(http.StatusTooManyRequests, "Rate limit exceeded", nil)
)
