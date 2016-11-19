package api

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// InitOvrstatAPI binds all required handlers to the router and starts
// the ovrstat API
func InitOvrstatAPI(segmentAPIKey, port string) {
	// Bindings dependencies
	c := context.Background()
	r := mux.NewRouter().StrictSlash(true)

	// Binds the stats endpoints to the router
	InitStatsBindings(c, r, NewStatsService(segmentAPIKey))

	// Listen and serve on the port passed
	log.Fatal(http.ListenAndServe(port, r))
}
