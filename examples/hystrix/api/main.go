package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

func ping(c web.C, w http.ResponseWriter, r *http.Request) {
	resultChan := make(chan string, 1)

	errChan := hystrix.Go("call_backend", func() error {
		fmt.Println("call backend hystrix......")
		res, err := http.Get("http://localhost:8001/")

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
		nil,
		// func(err error) error {
		// 	fmt.Println("call backup......")
		// 	// do this when services are down
		// 	resultChan <- "Hello from backup function!"

		// 	return nil
		// },
	)

	// Block until we have a result or an error.
	select {
	case result := <-resultChan:
		log.Println("success:", result)
		fmt.Fprint(w, string(result))
		w.WriteHeader(http.StatusOK)
	case err := <-errChan:
		log.Println("failure:", err)
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}

func main() {
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
