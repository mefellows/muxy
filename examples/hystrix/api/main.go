package main

import (
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/quipo/statsd"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
	"net/http"
	"os"
)

//var s *statsd.StatsdBuffer
var s *statsd.StatsdClient

func ping(c web.C, w http.ResponseWriter, r *http.Request) {
	hystrix.Go("call_backend", func() error {
		res, err := http.Get("http://backend/")
		if err == nil && res != nil {
			s.Incr("ok", 1)
			res.Write(w)
			return nil
		}
		s.Incr("err", 1)
		return err
	}, func(err error) error {
		// do this when services are down
		s.Incr("backup", 1)
		fmt.Printf("Hello from backup function!\n")
		fmt.Fprintf(w, "Hello from backup function!")
		return nil
	})
}

func main() {
	fmt.Println("Connecting to Statsd server...")

	prefix := "muxy."
	statsdclient := statsd.NewStatsdClient("statsd:8125", prefix)
	s = statsdclient
	s.Incr("test", 100)
	s.Gauge("something", 500)
	statsdclient.CreateSocket()
	//interval := time.Second * 2 // aggregate stats and flush every 2 seconds
	//s = statsd.NewStatsdBuffer(interval, statsdclient)

	if s == nil {
		fmt.Println("Could not connect to statsd server, exiting")
		os.Exit(1)
	}

	hystrix.ConfigureCommand("call_backend", hystrix.CommandConfig{
		Timeout:               1500,
		MaxConcurrentRequests: 100,
		ErrorPercentThreshold: 25,
	})

	goji.Abandon(middleware.Logger)
	goji.Get("/*", ping)
	goji.Serve()
	fmt.Printf("Closing...")
	s.Close()
}
