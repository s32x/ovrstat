package main

import (
	"os"

	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/labstack/echo"
	"github.com/sdwolfe32/ovrstat/api"
	"github.com/sirupsen/logrus"
)

// Port is the port the server will run on
var Port = getEnv("PORT", "8080")

func main() {
	logger := logrus.New()
	logger.Level = logrus.DebugLevel
	logger.Formatter = new(logrus.JSONFormatter)
	logger.Info("Ovrstat 0.3 - Simple Overwatch Stats API")

	// Create a server Builder and bind the endpoints
	logger.Info("Binding HTTP endpoints to router")
	e := echo.New()
	o := api.New(logger)
	e.GET("/stats/pc/:area/:tag", o.Stats)
	e.GET("/stats/:area/:tag", o.Stats)

	// Listen on the specified port
	logger.WithField("port", Port).Info("Listening for requests...")
	e.Logger.Fatal(e.Start(":" + Port))
}

// getEnv retrieves variables from the environment and falls back
// to a passed fallback variable if it isn't already set
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
