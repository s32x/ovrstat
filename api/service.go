package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sdwolfe32/ovrstat/ovrstat"
	"github.com/sdwolfe32/slimhttp"
	"github.com/sirupsen/logrus"
)

// OvrstatService defines all functionality for serving
// Overwatch stats
type OvrstatService interface {
	PCStats(r *http.Request) (interface{}, error)
	ConsoleStats(r *http.Request) (interface{}, error)
}

// ovrstatService contains all required dependencies for an
// OvrstatService
type ovrstatService struct{ log *logrus.Entry }

// NewOvrstatService initializes a new OvrstatService
func NewOvrstatService(log *logrus.Logger) OvrstatService {
	return &ovrstatService{log: log.WithField("service", "ovrstat")}
}

// PCStats handles serving Overwatch stats data for PC
func (o *ovrstatService) PCStats(r *http.Request) (interface{}, error) {
	l := o.log.WithField("handler", "pc").WithField("ip_address", r.RemoteAddr)
	l.Debug("New PC Stats request received")

	// Extracts request variables
	l.Debug("Extracting request variables")
	v := mux.Vars(r)
	region, tag := v["region"], v["tag"]

	// Verifies request variables exist
	l.Debug("Verifying existence of required variables")
	if region == "" || tag == "" {
		l.Error("Missing required variables")
		return nil, slimhttp.ErrorBadRequest
	}
	l = l.WithField("region", region).WithField("tag", tag)

	// Performs a full stats lookup
	l.Debug("Performing Stats lookup")
	stats, err := ovrstat.PCStats(region, tag)
	if err != nil {
		if err == ovrstat.PlayerNotFound {
			l.WithError(err).Error("Player not found")
			return nil, slimhttp.ErrorNotFound
		}
		l.WithError(err).Error("An error occurred during lookup")
		return nil, slimhttp.NewError("An internal error occurred",
			http.StatusInternalServerError, err)
	}

	// Returns the successful stats lookup
	l.Info("Returning successful stats lookup")
	return stats, nil
}

// ConsoleStats handles serving Overwatch stats data for Console
func (o *ovrstatService) ConsoleStats(r *http.Request) (interface{}, error) {
	l := o.log.WithField("handler", "console").WithField("ip_address", r.RemoteAddr)
	l.Debug("New Console Stats request received")

	// Extracts request variables
	l.Debug("Extracting request variables")
	v := mux.Vars(r)
	platform, tag := v["platform"], v["tag"]

	// Verifies request variables exist
	l.Info("Verifying existence of required variables")
	if platform == "" || tag == "" {
		l.Error("Missing required variables")
		return nil, slimhttp.ErrorBadRequest
	}
	l = l.WithField("platform", platform).WithField("tag", tag)

	// Performs a full stats lookup
	l.Info("Performing Stats lookup")
	stats, err := ovrstat.ConsoleStats(platform, tag)
	if err != nil {
		if err == ovrstat.PlayerNotFound {
			l.WithError(err).Error("Player not found")
			return nil, slimhttp.ErrorNotFound
		}
		l.WithError(err).Error("An error occurred during lookup")
		return nil, slimhttp.NewError("An error occurred during lookup",
			http.StatusInternalServerError, err)
	}

	// Returns the successful stats lookup
	l.Info("Returning successful stats lookup")
	return stats, nil
}
