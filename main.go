package main

import (
	"os"

	"github.com/sdwolfe32/ovrstat/service"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	log.Formatter = new(logrus.JSONFormatter)
	log.Info("Ovrstat 0.2 - Simple Overwatch Stats API")

	// Retrieve the port from the environment
	log.Info("Retrieving PORT from environment")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Create a server Builder and bind the endpoints
	log.Info("Binding HTTP endpoints to router")
	b := service.NewBuilder(1, port)
	service.BindEndpoints(b, log)

	// Listen on the specified port
	log.Info("Listening and Serving on port " + port)
	log.Fatal(b.ListenAndServe())
}
