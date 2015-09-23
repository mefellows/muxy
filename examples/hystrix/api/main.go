package main

import (
	"fmt"
	"net/http"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

func ping(c web.C, w http.ResponseWriter, r *http.Request) {
	hystrix.Go("call_backend", func() error {
		res, err := http.Get("http://backend/")
		if err == nil && res != nil {
			fmt.Fprintln(w, "Response from backend: \n")
			fmt.Println("Response from backend: ")
			res.Write(w)
			return nil
		}
		return err
	}, func(err error) error {
		// do this when services are down
		fmt.Fprintf(w, "Hello from backup function!")
		return nil
	})
}

func main() {
	hystrix.ConfigureCommand("call_backend", hystrix.CommandConfig{
		Timeout:               1500,
		MaxConcurrentRequests: 100,
		ErrorPercentThreshold: 25,
	})

	goji.Get("/*", ping)
	goji.Serve()
}
