package main /* import "s32x.com/ovrstat" */

import (
	"os"

	"s32x.com/ovrstat/api"
)

func main() { api.Start(getenv("PORT", "8080")) }

// getenv retrieves a variable from the environment and falls back to a passed
// default value if the key doesn't exist
func getenv(key, def string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return def
}
