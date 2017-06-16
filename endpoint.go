package main

import (
	"encoding/json"
	"net/http"
)

// Endpoint is a service endpoint that recieves an http request and returns
// either a successfully populated JSON response body or a JSON encoded error.
type Endpoint func(*http.Request) (interface{}, error)

// endpointWrapper turns an Endpoint into a standard http.Handlerfunc
func endpointWrapper(e Endpoint) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := e(r)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err)
			return
		}
		// Writes a successful json response body
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(res)
	}
}
