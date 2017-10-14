package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// Router defines all functionality for our api service router
type Router interface {
	HandleEndpoint(pattern string, endpoint Endpoint) *mux.Route
	ListenAndServe(port string) error
}

// Router contains the router and all rate limiting details
type router struct{ *mux.Router }

// NewRouter generates a new Router that will be used to bind
// handlers to the *mux.Router
func NewRouter(version int) Router {
	return &router{mux.NewRouter().PathPrefix("/v" + strconv.Itoa(version)).Subrouter()}
}

// HandleEndpoint binds a new Endpoint handler to the router
func (r *router) HandleEndpoint(pattern string, endpoint Endpoint) *mux.Route {
	return r.HandleFunc(pattern, endpointWrapper(endpoint))
}

// ListenAndServe applies CORS headers and starts the server
// using the *mux.Router
func (r *router) ListenAndServe(port string) error {
	// Create the basic HTTP server with base parameters
	srv := &http.Server{
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	// Apply CORS headers
	srv.Handler = handlers.CORS(
		handlers.AllowedHeaders([]string{"Authorization", "Content-Type"}),
		handlers.AllowedMethods([]string{"GET", "OPTIONS", "HEAD"}),
		handlers.AllowedOrigins([]string{"*"}),
	)(r)

	// Set the port to run on and serve
	srv.Addr = ":" + port
	return srv.ListenAndServe()
}

// An Endpoint is a service endpoint that receives a request and returns either
// a successfully processed response-body or an Error. In either case both
// responses are encoded and returned to the user with the appropriate status
// code
type Endpoint func(*http.Request) (interface{}, error)

// endpointWrapper transforms an endpoint into a standard http.Handlerfunc
func endpointWrapper(e Endpoint) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		// Handle the request and respond appropriately
		res, err := e(req)
		if err != nil {
			if e, ok := err.(*Error); ok {
				encodeJSON(rw, e.StatusCode, e)
			} else {
				encodeJSON(rw, http.StatusInternalServerError,
					NewError("An error has occurred", http.StatusInternalServerError, err))
			}
			return
		}
		encodeJSON(rw, http.StatusOK, res)
	}
}

// encodeJSON encodes the response to JSON and writes it to the ResponseWriter
func encodeJSON(w http.ResponseWriter, status int, res interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(res)
}
