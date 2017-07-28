package service

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sdwolfe32/ovrstat/ovrstat"
	"github.com/sirupsen/logrus"
)

// Service contains all required context needed to serve Ovrstat requests
type Service struct{ log *logrus.Entry }

// BindEndpoints binds the endpoints needed to serve ovrstat data to the
// passed Builders mux.Router
func BindEndpoints(b *Builder, log *logrus.Logger) {
	s := &Service{log: log.WithField("service", "ovrstat")}
	b.HandleEndpoint("/stats/pc/{region}/{tag}", s.pcHandler)
	b.HandleEndpoint("/stats/{platform}/{tag}", s.consoleHandler)
}

// pcHandler handles serving Overwatch stats data for PC
func (s *Service) pcHandler(r *http.Request) (interface{}, error) {
	l := s.log.WithField("handler", "pc").WithField("ip_address", r.RemoteAddr)
	l.Info("New PC Stats request received")

	// Extracts request variables
	l.Info("Extracting request variables")
	v := mux.Vars(r)
	region, tag := v["region"], v["tag"]

	// Verifies request variables
	l.Info("Verifying existance of required variables")
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
		l.WithError(err).Error("An error occured during lookup")
		return nil, NewError(err.Error(), http.StatusInternalServerError)
	}

	// Returns the successful stats lookup
	l.Info("Returning successful stats lookup")
	return stats, nil
}

// consoleHandler handles serving Overwatch stats data for Console
func (s *Service) consoleHandler(r *http.Request) (interface{}, error) {
	l := s.log.WithField("handler", "console").WithField("ip_address", r.RemoteAddr)
	l.Info("New Console Stats request received")

	// Extracts request variables
	l.Info("Extracting request variables")
	v := mux.Vars(r)
	platform, tag := v["platform"], v["tag"]

	// Verifies request variables
	l.Info("Verifying existance of required variables")
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
		l.WithError(err).Error("An error occured during lookup")
		return nil, NewError(err.Error(), http.StatusInternalServerError)
	}

	// Returns the successful stats lookup
	l.Info("Returning successful stats lookup")
	return stats, nil
}
