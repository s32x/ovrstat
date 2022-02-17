package main

import (
	"log"
	"os"

	"github.com/s32x/ovrstat/service"
)

func main() {
	// Start a new service
	service.Start(getenv("PORT", "8080")) // The port the server will run on
}

// getenv attempts to retrieve and return a variable from the environment. If it
// fails it will either crash or failover to a passed default value
func getenv(key string, def ...string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	if len(def) == 0 {
		log.Fatalf("%s not defined in environment", key)
	}
	return def[0]
}
