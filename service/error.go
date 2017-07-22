package service

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var (
	ErrorBadRequest          = NewError("Bad Request", http.StatusBadRequest)
	ErrorNotFound            = NewError("Not Found", http.StatusNotFound)
	ErrorInternalServerError = NewError("Internal Server Error", http.StatusInternalServerError)
)

// Error is a structure used to express API errors in JSON
type Error struct {
	Message string `json:"error,omitempty"`
	Status  int    `json:"status,omitempty"`
}

// NewError is a generic http error type used for all error responses
func NewError(message string, status int) *Error {
	return &Error{Message: message, Status: status}
}

// Error satisfies the error interface by returning a string representation of
// the error
func (e *Error) Error() string {
	return fmt.Sprintf("Error %d: '%s'", e.Status, e.Message)
}

// RespondWithJSON writes the json representation of the Error to the
// ResponseWriter
func (e *Error) RespondWithJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(e.Status)
	json.NewEncoder(w).Encode(e)
}
