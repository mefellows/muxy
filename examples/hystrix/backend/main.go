package main

import (
	"fmt"
	"net/http"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
)

func ping(c web.C, w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Hello from real API\n")
	fmt.Fprintf(w, "Hello from real API")
}

func main() {
	goji.Abandon(middleware.Logger)
	goji.Get("/*", ping)
	goji.Serve()
}
