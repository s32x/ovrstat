package api

import (
	"context"

	analytics "github.com/segmentio/analytics-go"
)

// HealthcheckService
type HealthcheckService interface {
	GetHealthcheck(context.Context, interface{}) (interface{}, error)
}

// healthcheckService
type healthcheckService struct {
	analytics *analytics.Client
}

// NewStatsService generates a new StatsService endpoint
func NewHealthcheckService(a *analytics.Client) HealthcheckService {
	return healthcheckService{analytics: a}
}

// getHealthcheckResponse
type getHealthcheckResponse struct {
	Status string `json:"status"`
}

// GetHealthcheck
func (s healthcheckService) GetHealthcheck(_ context.Context, req interface{}) (interface{}, error) {
	return &getHealthcheckResponse{Status: "OK"}, nil
}
