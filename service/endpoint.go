package service

import (
	"encoding/json"
	"net/http"
)

// Endpoint is a service endpoint that receives an http request and returns
// either a successfully populated JSON response body or a JSON encoded error
type Endpoint func(*http.Request) (interface{}, error)

// endpointWrapper transforms an Endpoint into a standard http.Handlerfunc
func endpointWrapper(e Endpoint) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := e(r)
		if err != nil {
			if e, ok := err.(*Error); ok {
				e.RespondWithJSON(w)
			} else {
				ErrorInternalServerError.RespondWithJSON(w)
			}
			return
		}

		// Writes a successful json response body to the ResponseWriter
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
	}
}
