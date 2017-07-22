package service

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

// Builder is a type used for building out service endpoints
type Builder struct {
	router *mux.Router
	port   string
}

// NewBuilder generates a new Builder that will be used to bind endpoints and
// launch an HTTP server
func NewBuilder(version int, port string) *Builder {
	return &Builder{
		router: mux.NewRouter().PathPrefix("/v" + strconv.Itoa(version)).Subrouter(),
		port:   port,
	}
}

// HandleEndpoint binds a new Endpoint to the router
func (s *Builder) HandleEndpoint(path string, endpoint Endpoint) *mux.Route {
	return s.router.HandleFunc(path, endpointWrapper(endpoint))
}

// ListenAndServe starts the server using the non-exported router. It applies
// simple CORS requirements and launches a basic HTTP server
func (s *Builder) ListenAndServe() error {
	// Create the basic HTTP server with base parameters
	srv := &http.Server{
		Handler:      s.router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	// Apply CORS headers
	srv.Handler = handlers.CORS(
		handlers.AllowedHeaders([]string{"Content-Type"}),
		handlers.AllowedMethods([]string{"OPTIONS", "GET", "HEAD"}),
		handlers.AllowedOrigins([]string{"*"}), // Allow all origins
	)(s.router)

	// Set the port and listen for requests
	srv.Addr = ":" + s.port
	return srv.ListenAndServe()
}
