package muxy

import (
	"fmt"
	s "github.com/mefellows/muxy/symptom"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
)

// Listen on a given port/ip
type Config struct{}

type Muxy struct {
	c        Config
	symptoms []Symptom
}

func New() *Muxy {
	return &Muxy{}
}

func (m *Muxy) Run() {
	log.Println("Running Muxy server on port 1234")

	// Read in test fixture settings as state into the Muxy server

	// Load all plugins

	// -> Symptoms
	// Read in symptoms + symptom configs (Use Hashi's config type or custom DSL?)
	symptom := &s.ShittyNetworkSymptom{}

	// -> Middlewares (stats etc.)

	// Interrupt handler
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	// Startup proxy
	go func() {
		mux := http.NewServeMux()
		symptom.Setup()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			director := func(req *http.Request) {
				req = r
				req.URL.Scheme = "http"
				req.URL.Host = "127.0.0.1:8282"
			}
			log.Println("request received")
			symptom.Muck()

			proxy := &httputil.ReverseProxy{Director: director}
			proxy.ServeHTTP(w, r)
		})
		err := http.ListenAndServe(":1234", mux)
		if err != nil {
			log.Println(fmt.Sprintf("ListenAndServe error: ", err))
		}
	}()

	// Block until a signal is received.
	<-c
	log.Println("Shutting down Muxy...")

	symptom.Teardown()
}
