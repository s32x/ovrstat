package service /* import "s32x.com/ovrstat/service" */

import (
	"net/http"
	"strings"

	packr "github.com/gobuffalo/packr/v2"
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
func Start(port, env string) {
	// Create a new echo Echo and bind all middleware
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.Gzip())

	// Perform HTTP redirects and serve the web index if being hosted in prod
	if strings.Contains(strings.ToLower(env), "prod") {
		e.Pre(middleware.HTTPSNonWWWRedirect())

		// Serve the static web content
		wb := packr.New("web box", "./web")
		e.GET("*", echo.WrapHandler(http.FileServer(wb)))
	}

	// Handle stats API requests
	e.GET("/stats/pc/:area/:tag", stats)
	e.GET("/stats/:area/:tag", stats)

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
