package main

import (
	"fmt"
	"net/http"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

func ping(c web.C, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from real API")
}

func main() {
	goji.Get("/*", ping)
	goji.Serve()
}
