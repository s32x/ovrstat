package api

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	analytics "github.com/segmentio/analytics-go"
)

// InitOvrstatAPI binds all required handlers to the router and starts
// the ovrstat API
func InitOvrstatAPI(segmentAPIKey, port string) {
	// Bindings dependencies
	c := context.Background()
	r := mux.NewRouter().StrictSlash(true)

	// Get segment API key for tracking API hits
	var a *analytics.Client
	if segmentAPIKey != "" {
		a = analytics.New(segmentAPIKey)
	}

	// Binds the stats endpoints to the router
	InitStatsBindings(c, r, NewStatsService(a))
	InitHealthcheckBindings(c, r, NewHealthcheckService(a))

	// Listen and serve on the port passed
	log.Fatal(http.ListenAndServe(port, r))
}
