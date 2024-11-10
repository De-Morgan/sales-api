package response

import (
	"errors"
	"math"
)

type PageDocument[T any] struct {
	Items      []T `json:"items"`
	Total      int `json:"total_items"`
	TotalPages int `json:"total_pages"`
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
}

func NewPageDocument[T any](items []T, total, page, pageSize int) PageDocument[T] {
	totalPage := math.Ceil(float64(total) / float64(pageSize))
	return PageDocument[T]{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: int(totalPage),
	}
}

type Success[T any] struct {
	Status bool `json:"status"`
	Data   T    `json:"data"`
}

func NewSuccess[T any](data T) Success[T] {
	return Success[T]{
		Status: true,
		Data:   data,
	}
}

// ErrorDocument is the form used for API responses from failures in the API.
type ErrorDocument struct {
	Error  string            `json:"error"`
	Fields map[string]string `json:"fields,omitempty"`
}

// Error is used to pass an error during the request through the
// application with web specific context.
type Error struct {
	Err    error
	Status int
}

// NewError wraps a provided error with an HTTP status code. This
// function should be used when handlers encounter expected errors.
func NewError(err error, status int) error {
	return &Error{
		err, status,
	}
}

// Error implements the error interface. It uses the default message of the
// wrapped error. This is what will be shown in the services' logs.
func (e *Error) Error() string {
	return e.Err.Error()
}

// IsError checks if an error of type Error exists.
func IsError(err error) bool {
	var re *Error
	return errors.As(err, &re)
}

// GetError returns a copy of the Error pointer.
func GetError(err error) *Error {
	var re *Error
	if !errors.As(err, &re) {
		return nil
	}
	return re
}
