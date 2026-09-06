package domain

import (
	"errors"
	"net/http"
)

type AppError struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewNotFoundError(msg string) *AppError {
	return &AppError{
		Code:    http.StatusNotFound,
		Message: msg,
		Err:     errors.New(msg),
	}
}

func NewBadRequestError(msg string) *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: msg,
		Err:     errors.New(msg),
	}
}

func NewInternalError(err error) *AppError {
	return &AppError{
		Code:    http.StatusInternalServerError,
		Message: "internal error",
		Err:     err,
	}
}
