package api /* import "s32x.com/ovrstat/api" */

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"s32x.com/ovrstat/ovrstat"
)

var (
	// ErrPlayerNotFound is thrown when a request is made for a player that doesn't exist
	ErrPlayerNotFound = echo.NewHTTPError(http.StatusNotFound, "Player not found")
	// ErrFailedLookup is thrown when there is an error retrieving an accounts stats
	ErrFailedLookup = echo.NewHTTPError(http.StatusInternalServerError, "Failed to perform lookup")
)

// Start starts the ovrstat API service using the passed params
func Start(port string) {
	// Create a server Builder and bind the endpoints
	e := echo.New()
	e.HideBanner = true
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
