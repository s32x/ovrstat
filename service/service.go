package service /* import "s32x.com/ovrstat/service" */

import (
	"embed"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//go:embed static/*
var static embed.FS

// Start starts serving the service on the passed port
func Start(port string) {
	// Create a new echo Echo and bind all middleware
	e := echo.New()
	e.HideBanner = true

	// Bind middleware
	e.Pre(middleware.RemoveTrailingSlashWithConfig(
		middleware.TrailingSlashConfig{
			RedirectCode: http.StatusPermanentRedirect,
		}))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Pre(middleware.Secure())
	e.Use(middleware.Gzip())
	e.Use(middleware.CORS())

	// Serve the static web content on the base echo instance
	e.GET("/*", echo.WrapHandler(http.FileServer(http.FS(static))),
		middleware.Rewrite(map[string]string{"/*": "/static/$1"}))

	// Handle stats API requests
	e.GET("/stats/:platform/:tag", stats)
	e.GET("/healthcheck", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	// Listen on the specified port
	e.Logger.Fatal(e.Start(":" + port))
}
