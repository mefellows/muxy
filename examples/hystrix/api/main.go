package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/quipo/statsd"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

var s *statsd.StatsdClient

func ping(c web.C, w http.ResponseWriter, r *http.Request) {
	resultChan := make(chan string, 1)

	errChan := hystrix.Go("call_backend", func() error {
		res, err := http.Get(os.Getenv("API_HOST"))

		if err == nil && res != nil {
			if res.StatusCode > 400 {
				return fmt.Errorf("Error: %d", res.StatusCode)
			}
			body, err := ioutil.ReadAll(res.Body)

			if err == nil {
				resultChan <- string(body)
			}
		}

		return err
	},
		func(err error) error {
			resultChan <- "Call from backup function"
			return nil
		},
	)

	// Block until we have a result or an error.
	select {
	case result := <-resultChan:
		s.Incr("ok", 1)
		fmt.Fprint(w, string(result))
		w.WriteHeader(http.StatusOK)
	case <-errChan:
		s.Incr("err", 1)
		w.WriteHeader(http.StatusServiceUnavailable)
	}
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

	// Setup hystrix streams
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
	go http.ListenAndServe(net.JoinHostPort("", "8181"), hystrixStreamHandler)

	hystrix.ConfigureCommand("call_backend", hystrix.CommandConfig{
		Timeout:               25,
		MaxConcurrentRequests: 100,
		ErrorPercentThreshold: 25,
	})

	goji.Get("/*", ping)
	goji.Serve()
}
