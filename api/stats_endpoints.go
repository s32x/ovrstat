package api

import (
	"context"
	"errors"

	"github.com/sdwolfe32/ovrstat/goow"
	analytics "github.com/segmentio/analytics-go"
)

// StatsService
type StatsService interface {
	GetStats(context.Context, getStatsRequest) (*getStatsResponse, error)
}

// statsService
type statsService struct {
	sc *analytics.Client
}

// NewStatsService generates a new StatsService endpoint
func NewStatsService(segmentAPIKey string) StatsService {
	// Get segment API key for tracking API hits
	var client *analytics.Client
	if segmentAPIKey != "" {
		client = analytics.New(segmentAPIKey)
	}
	return statsService{
		sc: client,
	}
}

// getStatsRequest
type getStatsRequest struct {
	platform string
	region   string
	tag      string
}

// getStatsResponse
type getStatsResponse struct {
	*goow.PlayerStats
}

// GetStats
func (s statsService) GetStats(_ context.Context, req getStatsRequest) (*getStatsResponse, error) {
	// Check for required fields before
	if req.platform == "" || req.region == "" || req.tag == "" {
		return nil, errors.New("Required fields are missing")
	}

	if s.sc != nil {
		// Track stats lookup event
		s.sc.Track(&analytics.Track{
			Event:       "Player Stats Lookup",
			AnonymousId: "ovrstat",
			Properties: map[string]interface{}{
				"platform": req.platform,
				"region":   req.region,
				"tag":      req.tag,
			},
		})
	}

	// Get the players stats from blizzard
	stats, err := goow.GetPlayerStats(req.platform, req.region, req.tag)
	if err != nil {
		return nil, errors.New("There was an error retrieving stats")
	}

	// If the player doesn't exist let the user know
	if stats.Name == "" && stats.Level == 0 {
		return nil, errors.New("The requested player was not found")
	}

	// Return the stats data embedded in our getStatsResponse wrapper
	return &getStatsResponse{stats}, nil
}
