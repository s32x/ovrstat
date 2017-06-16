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
	r := mux.NewRouter()
	r.Handle(prefix+"/pc/{region}/{tag}", endpointWrapper(pcHandler))
	r.Handle(prefix+"/{platform}/{tag}", endpointWrapper(consoleHandler))

	// Spin up a TLS goroutine if a cert and key are found
	certFile, keyFile := os.Getenv("CERT_FILE"), os.Getenv("KEY_FILE")
	if certFile != "" && keyFile != "" {
		log.Println("Listening on port :443")
		go http.ListenAndServeTLS(":443", certFile, keyFile, r)
	}

	// Always run a basic http server
	log.Println("Listening on port :80")
	http.ListenAndServe(":80", r)
}

func pcHandler(r *http.Request) (interface{}, error) {
	v := mux.Vars(r)
	return ovrstat.PCStats(v["region"], v["tag"])
}

func consoleHandler(r *http.Request) (interface{}, error) {
	v := mux.Vars(r)
	return ovrstat.ConsoleStats(v["platform"], v["tag"])
}
