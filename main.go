package main

import (
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/sdwolfe32/ovrstat/api"
	"github.com/sdwolfe32/slimhttp"
)

func main() {
	logger := logrus.New()
	logger.Level = logrus.DebugLevel
	logger.Formatter = new(logrus.JSONFormatter)
	logger.Info("Ovrstat 0.3 - Simple Overwatch Stats API")

	// Create a server Builder and bind the endpoints
	logger.Info("Binding HTTP endpoints to router")
	r := slimhttp.NewRouter()
	o := api.NewOvrstatService(logger)
	r.HandleJSONEndpoint("/stats/pc/{region}/{tag}", o.PCStats).Methods(http.MethodGet)
	r.HandleJSONEndpoint("/stats/{platform}/{tag}", o.ConsoleStats).Methods(http.MethodGet)

	// Listen on the specified port
	port := getEnv("PORT", "8000")
	logger.Info("Listening and Serving on port " + port)
	logger.Fatal(r.ListenAndServe(port))
}

// getEnv retrieves variables from the environment and falls back
// to a passed fallback variable if it isn't already set
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
