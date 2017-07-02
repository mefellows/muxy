package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/quipo/statsd"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

var s *statsd.StatsdClient

func ping(c web.C, w http.ResponseWriter, r *http.Request) {
	s.Incr("backup", 1)
	fmt.Fprintf(w, "Hello from real API")
	w.WriteHeader(http.StatusOK)
}

func main() {
	// Setup Statsd client
	fmt.Println("Connecting to Statsd server...")

	prefix := "muxy."
	statsdclient := statsd.NewStatsdClient(os.Getenv("STATSD_HOST"), prefix)
	s = statsdclient
	statsdclient.CreateSocket()

	if s == nil {
		fmt.Println("Could not connect to statsd server, exiting")
		os.Exit(1)
	}

	goji.Get("/*", ping)
	goji.Serve()
}
