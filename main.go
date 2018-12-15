package main /* import "s32x.com/ovrstat" */

import (
	"log"
	"os"

	"s32x.com/ovrstat/service"
)

func main() {
	service.Start(
		getenv("PORT", "8080"),
		getenv("ENV", "dev"),
	)
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
