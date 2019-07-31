package main

import (
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"s32x.com/ovrstat/ovrstat"
)

var port = getenv("PORT", "8080")

func main() {
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
	e.Static("*", "./static")

	// Handle stats API requests
	e.GET("/stats/:platform/:tag", stats)
	e.GET("/healthcheck", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	// Listen on the specified port
	e.Logger.Fatal(e.Start(":" + port))
}

// stats handles retrieving and serving Overwatch stats in JSON
func stats(c echo.Context) error {
	// Perform a full player stats lookup
	stats, err := ovrstat.Stats(c.Param("platform"), c.Param("tag"))
	if err != nil {
		if err == ovrstat.ErrPlayerNotFound {
			return newErr(http.StatusNotFound, "Player not found")
		}
		return newErr(http.StatusInternalServerError,
			errors.Wrap(err, "Failed to retrieve player stats"))
	}
	return c.JSON(http.StatusOK, stats)
}

// newErr creates and returns a new echo HTTPError with the passed status code
// and optional message. Message expected to be of type string or error
func newErr(code int, message ...interface{}) error {
	if len(message) > 0 {
		switch v := message[0].(type) {
		case error:
			return echo.NewHTTPError(code, v.Error())
		case string:
			return echo.NewHTTPError(code, v)
		}
	}
	return echo.NewHTTPError(code, "An error has occurred")
}

// getenv attempts to retrieve and return a variable from the environment. If it
// fails it will either crash or failover to a passed default value
func getenv(key string, def ...string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	if len(def) == 0 {
		log.Fatalf("%s not defined in environment", key)
	}
	return def[0]
}
