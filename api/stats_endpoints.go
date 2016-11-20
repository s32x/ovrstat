package api

import (
	"context"
	"errors"
	"time"

	cache "github.com/patrickmn/go-cache"
	"github.com/sdwolfe32/ovrstat/goow"
	analytics "github.com/segmentio/analytics-go"
)

// StatsService
type StatsService interface {
	GetStats(context.Context, *getStatsRequest) (interface{}, error)
}

// statsService
type statsService struct {
	analytics *analytics.Client
	cache     *cache.Cache
}

// NewStatsService generates a new StatsService endpoint
func NewStatsService(a *analytics.Client) StatsService {
	return statsService{
		analytics: a,
		// Retains cached stats for a minute. Flushes every 10 seconds
		cache: cache.New(time.Minute, 10*time.Second),
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
func (s statsService) GetStats(_ context.Context, req *getStatsRequest) (interface{}, error) {
	// Check for required fields before
	if req.platform == "" || req.region == "" || req.tag == "" {
		return nil, errors.New("Required fields are missing")
	}

	if s.analytics != nil {
		// Track stats lookup event
		s.analytics.Track(&analytics.Track{
			Event:       "Player Stats Lookup",
			AnonymousId: "ovrstat",
			Properties: map[string]interface{}{
				"platform": req.platform,
				"region":   req.region,
				"tag":      req.tag,
			},
		})
	}

	// Attempts to get the players stats from our in-memory cache
	if x, found := s.cache.Get(req.platform + "/" + req.region + "/" + req.tag); found {
		stats := x.(*goow.PlayerStats)
		return &getStatsResponse{stats}, nil
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

	// Stores the players stats to our in-memory cache
	s.cache.Set(req.platform+"/"+req.region+"/"+req.tag, stats, cache.DefaultExpiration)

	// Return the stats data embedded in our getStatsResponse wrapper
	return &getStatsResponse{stats}, nil
}
