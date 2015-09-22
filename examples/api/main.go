package main

import (
	"fmt"
	"net/http"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	// "github.com/zenazn/goji/web/middleware"
)

func ping(c web.C, w http.ResponseWriter, r *http.Request) {
	// hystrix.ConfigureCommand("call_backend", hystrix.CommandConfig{
	// 	Timeout: 500,
	// })

	hystrix.Go("call_backend", func() error {
		res, err := http.Get("http://backend/")
		if err != nil {
			fmt.Fprintln(w, "Response from backend: \n")
			res.Write(w)
		}
		return err
	}, func(err error) error {
		// do this when services are down
		fmt.Fprintf(w, "Hello from backup function!")
		return nil
	})
}

func main() {
	goji.Get("/*", ping)
	// goji.Abandon(middleware.Logger)
	goji.Serve()
}
