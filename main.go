package main

import (
	"github.com/kataras/iris"
	"github.com/sdwolfe32/ovrstat/goow"
)

func main() {
	iris.Get("/v1/:platform/:region/:tag", func(c *iris.Context) {
		platform := c.Param("platform")
		region := c.Param("region")
		tag := c.Param("tag")

		if platform == "" || region == "" || tag == "" {
			c.JSON(iris.StatusNotFound, iris.Map{"Error": "Required fields are missing"})
			return
		}

		player, err := goow.GetPlayerStats(platform, region, tag)
		if err != nil {
			c.JSON(iris.StatusNotFound, iris.Map{"Error": "There was an error retrieving stats"})
			return
		}

		if player.Name == "" && player.Level == 0 {
			c.JSON(iris.StatusNotFound, iris.Map{"Error": "The requested player was not found"})
			return
		}

		c.JSON(iris.StatusOK, player)
	})
	iris.Listen("0.0.0.0:7000")
}
