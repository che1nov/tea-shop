package errors

import "fmt"

type ErrorCode string

const (
	ErrInvalidInput       ErrorCode = "INVALID_INPUT"
	ErrNotFound           ErrorCode = "NOT_FOUND"
	ErrUnauthorized       ErrorCode = "UNAUTHORIZED"
	ErrInternal           ErrorCode = "INTERNAL_ERROR"
	ErrConflict           ErrorCode = "CONFLICT"
	ErrServiceUnavailable ErrorCode = "SERVICE_UNAVAILABLE"
)

type AppError struct {
	Code    ErrorCode
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func New(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

func NewWithErr(code ErrorCode, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
