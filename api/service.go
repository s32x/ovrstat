package api

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/sdwolfe32/ovrstat/ovrstat"
	"github.com/sirupsen/logrus"
)

var (
	// ErrPlayerNotFound is thrown when a request is made for a player that doesn't exist
	ErrPlayerNotFound = echo.NewHTTPError(http.StatusNotFound, "Player not found")

	// ErrFailedLookup is thrown when there is an error retrieving an accounts stats
	ErrFailedLookup = echo.NewHTTPError(http.StatusInternalServerError, "Failed to perform lookup")
)

// OvrstatService contains all required dependencies for performing
// Overwatch stats lookups
type OvrstatService struct{ log *logrus.Entry }

// New generates and returns a new ovrstatService reference
func New(log *logrus.Logger) *OvrstatService {
	return &OvrstatService{log: log.WithField("service", "ovrstat")}
}

// Stats handles serving Overwatch stats data for PC
func (o *OvrstatService) Stats(c echo.Context) error {
	l := o.log.WithField("handler", "pc").WithField("ip_address", c.RealIP())
	l.Debug("New PC Stats request received")

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

	// Returns the successful stats lookup
	l.Info("Returning successful stats lookup")
	return c.JSON(http.StatusOK, stats)
}
