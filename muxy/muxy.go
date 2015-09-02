package muxy

import (
	"fmt"
	s "github.com/mefellows/muxy/symptom"
  "github.com/mefellows/muxy/config"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
)

const (
	DEFAULT_PORT = 8123
	DEFAULT_HOST = "0.0.0.0"
)

type MuxyConfig struct {
	Port          int    // Which port to listen on
	Host          string // Which network host/ip to listen on
	ProxyPort     int    // Which port proxy
	ProxyHost     string // Which network host/ip to proxy
	ProxyProtocol string // Which protocol to proxy (http, tcp, ...)
	Insecure      bool   // Enable/Disable TLS between Muxy <-> Proxied Host
	Ssl           bool   // Enable/Disable TLS between client <-> Muxy
	RawConfig     *config.RawConfig
	ConfigFile    string // Path to YAML Configuration File
}

type Muxy struct {
	config *MuxyConfig

	symptoms    []Symptom
	middlewares []Middleware
}

func New(config *MuxyConfig) *Muxy {
	return &Muxy{config: config}
}

func NewWithDefaultMuxyConfig() *Muxy {
	c := &MuxyConfig{
		Port: DEFAULT_PORT,
		Host: DEFAULT_HOST,
	}
	return &Muxy{config: c}
}

func (m *Muxy) Run() {
	log.Println(fmt.Sprintf("Running Muxy server on port %d", m.config.Port))

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
				req.URL.Scheme = m.config.ProxyProtocol
				req.URL.Host = fmt.Sprintf("%s:%d", m.config.ProxyHost, m.config.ProxyPort)
				log.Println("Proxying : ", fmt.Sprintf("%s://%s:%d", m.config.ProxyProtocol, m.config.ProxyHost, m.config.ProxyPort))
			}
			log.Println("request received")
			symptom.Muck()

			proxy := &httputil.ReverseProxy{Director: director}
			proxy.ServeHTTP(w, r)
		})
		err := http.ListenAndServe(fmt.Sprintf("%s:%d", m.config.Host, m.config.Port), mux)
		if err != nil {
			log.Println(fmt.Sprintf("ListenAndServe error: ", err))
		}
	}()

	// Block until a signal is received.
	<-c
	log.Println("Shutting down Muxy...")

	symptom.Teardown()
}
