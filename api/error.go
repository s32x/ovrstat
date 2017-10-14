package api

import (
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
)

var (
	ErrorBadRequest          = NewError("Bad Request", http.StatusBadRequest, nil)
	ErrorNotFound            = NewError("Not Found", http.StatusNotFound, nil)
	ErrorInternalServerError = NewError("Internal Server Error", http.StatusInternalServerError, nil)
)

// Error is a generic API error
type Error struct {
	Message    string `json:"message,omitempty"`
	StatusCode int    `json:"status,omitempty"`
	Err        string `json:"err,omitempty"`
}

// NewError is a generic http error type used for all error responses
func NewError(message string, statusCode int, err error) *Error {
	var errString string
	if err != nil {
		errString = err.Error()
	}
	return &Error{
		Message:    message,
		StatusCode: statusCode,
		Err:        errString,
	}
}

// Error returns a string representation of the error and
// helps to satisfy the error interface
func (e *Error) Error() string {
	return fmt.Sprintf("Error %d: '%s'", e.StatusCode, e.Message)
}

// Log will log the Error before returning it
func (e *Error) Log(log *logrus.Entry) *Error {
	log.WithFields(logrus.Fields{"error": e.Err, "status": e.StatusCode}).Error(e.Message)
	return e
}
