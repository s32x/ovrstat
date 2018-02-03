package main

import (
	"net/http"
	"os"
	"strconv"

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
	h := slimhttp.NewHealthcheckService(logger, "ovrstat.com")
	r.HandleJSONEndpoint("/stats/pc/{region}/{tag}", o.PCStats).Methods(http.MethodGet)
	r.HandleJSONEndpoint("/stats/{platform}/{tag}", o.ConsoleStats).Methods(http.MethodGet)
	r.HandleJSONEndpoint("/healthcheck", h.Healthcheck).Methods(http.MethodGet)

	// Listen on the specified port
	port, _ := strconv.Atoi(getEnv("PORT", "8000"))
	logger.Infof("Listening and Serving on port %v", port)
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
