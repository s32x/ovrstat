package api

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/sdwolfe32/ovrstat/ovrstat"
)

// Ovrstater defines all functionality for serving Overwatch stats
type Ovrstater interface {
	PCStats(r *http.Request) (interface{}, error)
	ConsoleStats(r *http.Request) (interface{}, error)
}

// Service contains all required context needed to serve Ovrstat requests
type ovrstater struct{ log *logrus.Entry }

// NewOvrstater initializes a new Ovrstater interface
func NewOvrstater(log *logrus.Logger) Ovrstater {
	return &ovrstater{log: log.WithField("service", "ovrstat")}
}

// PCStats handles serving Overwatch stats data for PC
func (o *ovrstater) PCStats(r *http.Request) (interface{}, error) {
	l := o.log.WithField("handler", "pc").WithField("ip_address", r.RemoteAddr)
	l.Info("New PC Stats request received")

	// Extracts request variables
	l.Info("Extracting request variables")
	v := mux.Vars(r)
	region, tag := v["region"], v["tag"]

	// Verifies request variables exist
	l.Info("Verifying existence of required variables")
	if region == "" || tag == "" {
		l.Error("Missing required variables")
		return nil, ErrorBadRequest
	}
	l = l.WithField("region", region).WithField("tag", tag)

	// Performs a full stats lookup
	l.Info("Performing Stats lookup")
	stats, err := ovrstat.PCStats(region, tag)
	if err != nil {
		if err == ovrstat.PlayerNotFound {
			l.WithError(err).Error("Player not found")
			return nil, ErrorNotFound
		}
		l.WithError(err).Error("An error occurred during lookup")
		return nil, NewError("An internal error occurred",
			http.StatusInternalServerError, err)
	}

	// Returns the successful stats lookup
	l.Info("Returning successful stats lookup")
	return stats, nil
}

// ConsoleStats handles serving Overwatch stats data for Console
func (o *ovrstater) ConsoleStats(r *http.Request) (interface{}, error) {
	l := o.log.WithField("handler", "console").WithField("ip_address", r.RemoteAddr)
	l.Info("New Console Stats request received")

	// Extracts request variables
	l.Info("Extracting request variables")
	v := mux.Vars(r)
	platform, tag := v["platform"], v["tag"]

	// Verifies request variables exist
	l.Info("Verifying existence of required variables")
	if platform == "" || tag == "" {
		l.Error("Missing required variables")
		return nil, ErrorBadRequest
	}
	l = l.WithField("platform", platform).WithField("tag", tag)

	// Performs a full stats lookup
	l.Info("Performing Stats lookup")
	stats, err := ovrstat.ConsoleStats(platform, tag)
	if err != nil {
		if err == ovrstat.PlayerNotFound {
			l.WithError(err).Error("Player not found")
			return nil, ErrorNotFound
		}
		l.WithError(err).Error("An error occurred during lookup")
		return nil, NewError("An error occurred during lookup",
			http.StatusInternalServerError, err)
	}

	// Returns the successful stats lookup
	l.Info("Returning successful stats lookup")
	return stats, nil
}
