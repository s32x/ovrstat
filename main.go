package main

import (
	"os"

	"github.com/sdwolfe32/ovrstat/api"
)

func main() {
	api.InitOvrstatAPI(os.Getenv("SEGMENT_API_KEY"), ":7000")
}
