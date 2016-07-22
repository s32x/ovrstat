package main

import (
	"os"

	"github.com/kataras/iris"
	"github.com/sdwolfe32/ovrstat/goow"
	"github.com/segmentio/analytics-go"
)

func main() {
	// Get segment API key for tracking API hits
	sc := new(analytics.Client)
	segmentAPIKey := os.Getenv("SEGMENT_API_KEY")
	if segmentAPIKey != "" {
		sc = analytics.New(segmentAPIKey)
	}

	iris.Get("/", func(c *iris.Context) {
		c.JSON(iris.StatusOK, iris.Map{"endpoints": "/v1/stats/{platform}/{region}/{tag}"})
	})
	iris.Get("/v1/stats/:platform/:region/:tag", func(c *iris.Context) {
		platform := c.Param("platform")
		region := c.Param("region")
		tag := c.Param("tag")

		// Check for required fields before
		if platform == "" || region == "" || tag == "" {
			c.JSON(iris.StatusNotFound, iris.Map{"Error": "Required fields are missing"})
			return
		}

		if segmentAPIKey != "" {
			// Track stats lookup event
			sc.Track(&analytics.Track{
				Event:       "Player Stats Lookup",
				AnonymousId: "ovrstat",
				Properties: map[string]interface{}{
					"platform": platform,
					"region":   region,
					"tag":      tag,
					"source":   c.RemoteAddr(),
				},
			})
		}

		// Get the players stats from blizzard
		player, err := goow.GetPlayerStats(platform, region, tag)
		if err != nil {
			c.JSON(iris.StatusNotFound, iris.Map{"Error": "There was an error retrieving stats"})
			return
		}

		// If the player doesn't exist let the user know
		if player.Name == "" && player.Level == 0 {
			c.JSON(iris.StatusNotFound, iris.Map{"Error": "The requested player was not found"})
			return
		}

		// Return the player struct with status ok
		c.JSON(iris.StatusOK, player)
	})
	iris.Listen("0.0.0.0:7000")
}
