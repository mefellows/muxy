package main

import (
	"fmt"
	"net/http"

	"github.com/quipo/statsd"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

var s *statsd.StatsdClient

func ping(c web.C, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from real API")
	w.WriteHeader(http.StatusOK)
}

func main() {
	goji.Get("/*", ping)
	goji.Serve()
}
