package api

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/sdwolfe32/ovrstat/ovrstat"
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
type Service struct{ log *logrus.Entry }

// New generates and returns a new ovrstatService reference
func New(log *logrus.Logger) *Service {
	return &Service{log.WithField("service", "ovrstat")}
}

// Overwatch handles serving Overwatch stats
func (o *Service) Overwatch(c echo.Context) error {
	l := o.log.WithField("handler", "overwatch")
	l.Debug("New Overwatch Stats request received")

	// Performs a full stats lookup
	l.Debug("Performing Stats lookup")
	stats, err := ovrstat.Stats(c.Param("area"), c.Param("tag"))
	if err != nil {
		if err == ovrstat.ErrPlayerNotFound {
			l.WithError(err).Error("Player not found")
			return ErrPlayerNotFound
		}
		l.WithError(err).Error("An error occurred during lookup")
		return ErrFailedLookup
	}

	// Returns the successful overwatch stats lookup
	l.Debug("Returning successful Overwatch Stats lookup")
	return c.JSON(http.StatusOK, stats)
}
