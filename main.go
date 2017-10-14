package main

import (
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/sdwolfe32/ovrstat/api"
)

func main() {
	log := logrus.New()
	log.Formatter = new(logrus.JSONFormatter)
	log.Info("Ovrstat 0.2 - Simple Overwatch Stats API")

	// Retrieve the port from the environment
	log.Info("Retrieving PORT from environment")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	// Create a server Builder and bind the endpoints
	log.Info("Binding HTTP endpoints to router")
	r := api.NewRouter(1)
	o := api.NewOvrstater(log)
	r.HandleEndpoint("/stats/pc/{region}/{tag}", o.PCStats).Methods(http.MethodGet)
	r.HandleEndpoint("/stats/{platform}/{tag}", o.ConsoleStats).Methods(http.MethodGet)

	// Listen on the specified port
	log.Info("Listening and Serving on port " + port)
	log.Fatal(r.ListenAndServe(port))
}
