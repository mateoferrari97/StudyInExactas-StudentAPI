package server

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type Error struct {
	StatusCode int    `json:"status"`
	Code       string `json:"code"`
	Message    string `json:"message"`
}

func NewError(message string, statusCode int) *Error {
	return NewErrorf(message, statusCode)
}

func NewErrorf(message string, statusCode int, args ...interface{}) *Error {
	return &Error{
		StatusCode: statusCode,
		Code:       strings.ReplaceAll(strings.ToLower(http.StatusText(statusCode)), " ", "_"),
		Message:    fmt.Sprintf(message, args...),
	}
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}

	if e.Message == "" {
		return e.Code
	}

	if e.Code == "" {
		return e.Message
	}

	return fmt.Sprintf("%d %s: %s", e.StatusCode, e.Code, e.Message)
}

func handleError(w http.ResponseWriter, err error) *Error {
	var hErr *Error
	if errors.As(err, &hErr) {
		return err.(*Error)
	}

	return NewErrorf(err.Error(), http.StatusInternalServerError)
}
