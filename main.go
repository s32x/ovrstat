package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/sdwolfe32/ovrstat/ovrstat"
)

const prefix = "/v1/stats"

func main() {
	log.Println("Ovrstat 0.1 - Simple Overwatch Stats API")

	// Retrieve the port from the environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	// Set up the pc/console endpoints
	r := mux.NewRouter()
	r.Handle(prefix+"/pc/{region}/{tag}", endpointWrapper(pcHandler))
	r.Handle(prefix+"/{platform}/{tag}", endpointWrapper(consoleHandler))

	// Always run a basic http server
	log.Printf("Listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func pcHandler(r *http.Request) (interface{}, error) {
	v := mux.Vars(r)
	return ovrstat.PCStats(v["region"], v["tag"])
}

func consoleHandler(r *http.Request) (interface{}, error) {
	v := mux.Vars(r)
	return ovrstat.ConsoleStats(v["platform"], v["tag"])
}
