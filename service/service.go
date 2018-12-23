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
	e.Pre(middleware.RemoveTrailingSlashWithConfig(
		middleware.TrailingSlashConfig{
			Skipper:      middleware.DefaultSkipper,
			RedirectCode: http.StatusPermanentRedirect,
		}))
	e.Pre(middleware.Secure())

	// Configure HTTP redirects and serve the web index if being hosted in prod
	if strings.Contains(strings.ToLower(env), "prod") {
		e.Pre(middleware.HTTPSNonWWWRedirect())
	}

	// Bind remaining middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())

	// Serve the static web content on the base echo instance
	wb := packr.New("web box", "./web")
	e.GET("*", echo.WrapHandler(http.FileServer(wb)))

	// Create the API group with separate middlewares
	api := e.Group("/stats")
	api.Use(middleware.CORS())

	// Handle stats API requests
	api.GET("/pc/:area/:tag", stats)
	api.GET("/:area/:tag", stats)

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
