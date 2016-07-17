package main

import (
	"github.com/kataras/iris"
)

func main() {
	iris.Get("/v1/:platform/:region/:tag", func(c *iris.Context) {
		platform := c.Param("platform")
		region := c.Param("region")
		tag := c.Param("tag")

		if platform == "" || region == "" || tag == "" {
			c.Error("Missing Required Fields", 400)
			return
		}

		player, err := goow.GetPlayerStats(platform, region, tag)
		if err != nil {
			c.Error("There was an error retrieving stats", 404)
		}

		c.JSON(iris.StatusOK, player)
	})
	iris.Listen(":8080")
}
