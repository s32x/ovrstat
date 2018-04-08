package api

import (
	"net/http"

	raven "github.com/getsentry/raven-go"
	"github.com/labstack/echo"
	tinystat "github.com/sdwolfe32/tinystat/client"
	"golang.org/x/sync/errgroup"
)

// Stats contains statistics on requests to the Ovrstat
// API
type Stats struct {
	DayLookups   int64 `json:"dayLookups"`
	MonthLookups int64 `json:"monthLookups"`
}

// Ovrstat handles serving Ovrstat stats
func (o *OvrstatService) Ovrstat(c echo.Context) error {
	l := o.log.WithField("handler", "stats")
	l.Debug("New Ovrstat Stats request received")

	// Retrieve stats from Tinystat
	var stats Stats
	var g errgroup.Group
	g.Go(func() (err error) {
		stats.DayLookups, err = tinystat.ActionCount("success", "24h")
		return
	})
	g.Go(func() (err error) {
		stats.MonthLookups, err = tinystat.ActionCount("success", "720h")
		return
	})
	if err := g.Wait(); err != nil {
		raven.CaptureError(err, nil)
		l.WithError(err).Error("An error occurred retrieving Tinystat stats")
	}

	// Returns the successful ovrstat Stats lookup
	l.Debug("Returning successful Ovrstat Stats lookup")
	return c.JSON(http.StatusOK, stats)
}
