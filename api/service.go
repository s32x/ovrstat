package api

import (
	"net/http"
	"strings"
	"time"

	raven "github.com/getsentry/raven-go"
	"github.com/labstack/echo"
	cache "github.com/patrickmn/go-cache"
	"github.com/sdwolfe32/ovrstat/ovrstat"
	tinystat "github.com/sdwolfe32/tinystat/client"
	"github.com/sirupsen/logrus"
)

// cacheExp specifies the length of time to hold on to
// player stats before retrieving new results
const cacheExp = time.Minute * 10

var (
	// ErrPlayerNotFound is thrown when a request is made for a player that doesn't exist
	ErrPlayerNotFound = echo.NewHTTPError(http.StatusNotFound, "Player not found")
	// ErrFailedLookup is thrown when there is an error retrieving an accounts stats
	ErrFailedLookup = echo.NewHTTPError(http.StatusInternalServerError, "Failed to perform lookup")
)

// Service contains all required dependencies for performing
// Overwatch stats lookups
type Service struct {
	log   *logrus.Entry
	cache *cache.Cache
}

// New generates and returns a new ovrstatService reference
func New(log *logrus.Logger) *Service {
	return &Service{
		log:   log.WithField("service", "ovrstat"),
		cache: cache.New(cacheExp, cacheExp),
	}
}

// Overwatch handles serving Overwatch stats
func (o *Service) Overwatch(c echo.Context) error {
	l := o.log.WithField("handler", "overwatch")
	l.Debug("New Overwatch Stats request received")

	// Generate a cache key and check the cache
	key := cacheKey(c.Param("area"), c.Param("tag"))
	if stats, ok := o.cache.Get(key); ok {
		// Returns the successful overwatch stats lookup
		l.Debug("Returning cached Overwatch Stats lookup")
		tinystat.CreateAction("success")
		return c.JSON(http.StatusOK, stats)
	}

	// Performs a full stats lookup
	l.Debug("Performing Stats lookup")
	stats, err := ovrstat.Stats(c.Param("area"), c.Param("tag"))
	if err != nil {
		tinystat.CreateAction("error")
		raven.CaptureError(err, nil)
		if err == ovrstat.ErrPlayerNotFound {
			l.WithError(err).Error("Player not found")
			return ErrPlayerNotFound
		}
		l.WithError(err).Error("An error occurred during lookup")
		return ErrFailedLookup
	}

	// Store the stats in cache for subsequent requests
	o.cache.SetDefault(key, stats)

	// Returns the successful overwatch stats lookup
	l.Debug("Returning successful Overwatch Stats lookup")
	tinystat.CreateAction("success")
	return c.JSON(http.StatusOK, stats)
}

// cacheKey takes an area and a tag and produces a
// unique cache key from the two
func cacheKey(area, tag string) string {
	return strings.Join([]string{area, tag}, "_")
}
