package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/starboy/ovrstat/ovrstat"
)

var (
	version = "0.4"
	port    = getEnv("PORT", "8080")

	// ErrPlayerNotFound is thrown when a request is made for a player that doesn't exist
	ErrPlayerNotFound = echo.NewHTTPError(http.StatusNotFound, "Player not found")

	// ErrFailedLookup is thrown when there is an error retrieving an accounts stats
	ErrFailedLookup = echo.NewHTTPError(http.StatusInternalServerError, "Failed to perform lookup")
)

func main() {
	fmt.Printf("Ovrstat %s - Simple Overwatch Stats API\n", version)

	// Create a server Builder and bind the endpoints
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Handle stats requests
	e.GET("/stats/pc/:area/:tag", stats)
	e.GET("/stats/:area/:tag", stats)

	// Serve the static web content
	e.Static("/", "web")
	e.Static("/assets", "web/assets")

	// Listen on the specified port
	fmt.Printf("Listening for requests on port : %s", port)
	e.Logger.Fatal(e.Start(":" + port))
}

// stats handles retrieving and serving Overwatch stats in JSON
func stats(c echo.Context) error {
	// Perform a full player stats lookup
	stats, err := ovrstat.Stats(c.Param("area"), c.Param("tag"))
	if err != nil {
		if err == ovrstat.ErrPlayerNotFound {
			return ErrPlayerNotFound
		}
		return ErrFailedLookup
	}
	return c.JSON(http.StatusOK, stats)
}

// getEnv retrieves variables from the environment and falls back to a passed
// fallback variable if it isn't already set
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
