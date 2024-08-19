package errs

import (
	"net/http"
)

type AppError struct {
	Code    int
	Message string
}

func (e AppError) Error() string {
	return e.Message
}

func NewNotFoundError(message string) error {

	return AppError{
		Code:    http.StatusNotFound,
		Message: message,
	}
}

func NewUnExpectedError() error {
	return AppError{
		Code:    http.StatusNotFound,
		Message: "Un Expected Error.",
	}
}

func NewBadRequestError(message string) error {
	return AppError{
		Code:    http.StatusBadRequest,
		Message: message,
	}
}
