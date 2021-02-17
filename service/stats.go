package service /* import "s32x.com/ovrstat/service" */

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"s32x.com/ovrstat/ovrstat"
)

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
