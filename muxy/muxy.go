package muxy

import (
	"fmt"
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

	// Load Configuration
	var c *config.Config
	var err error
	var confLoader *config.ConfigLoader
	if m.config.ConfigFile != "" {
		confLoader = &config.ConfigLoader{}
		c, err = confLoader.LoadFromFile(m.config.ConfigFile)
		if err != nil {
			log.Fatalf("Unable to read configuration file: %s", err.Error())
		}
	} else {
		log.Fatal("No config file provided")
	}

	// Load all plugins
	// TODO: Extract these plugin lifecycle events into the plugin lib?
	symptoms := make([]Symptom, len(c.Symptoms))
	for i, symptomConfig := range c.Symptoms {
		log.Printf("[DEBUG] Loading Symptom: %s", symptomConfig.Name)
		sf, ok := SymptomFactories.Lookup(symptomConfig.Name)

		if !ok {
			log.Fatalf("Unable to load symptom with name: %s", symptomConfig.Name)
		}

		//var s interface{}
		s, err := sf()
		if err != nil {
			log.Fatalf("Encountered error loading symptom: %v", err)
		}

		log.Printf("Symptom: %v\n", s)
		// apply config and validate
		err = confLoader.ApplyConfig(symptomConfig.Config, s)
		if err != nil {
			log.Fatalf("Encountered error applying configuration to symptom: %v", err)
		}
		err = confLoader.Validate(s)
		if err != nil {
			log.Fatalf("Encountered error validating symptom configuration: %v", err)
		}
		// TODO: Collapse this into the constructor, is it necessary??
		//s.Configure(&symptomConfig.Config)

		//symptoms[i] = s.(Symptom)
		symptoms[i] = s
	}

	// Setup all plugins...
	for _, s := range symptoms {
		s.Setup()
	}

	// -> Repeat for Middlewares (stats etc.)

	// Interrupt handler
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill)

	// Startup proxy
	// TODO: This should also be abstracted into a plugin
	go func() {
		mux := http.NewServeMux()
		//symptom.Setup()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			director := func(req *http.Request) {
				req = r
				req.URL.Scheme = m.config.ProxyProtocol
				req.URL.Host = fmt.Sprintf("%s:%d", m.config.ProxyHost, m.config.ProxyPort)
				log.Println("Proxying : ", fmt.Sprintf("%s://%s:%d", m.config.ProxyProtocol, m.config.ProxyHost, m.config.ProxyPort))
			}
			log.Println("request received")

			// Execute Middleware lifecycle hook: pre-muck
			//for _, s := range middlewares {
			//	go func() {
			//    s.Whatever(context...)
			//  }()
			//}

			//	symptom.Muck()
			for _, s := range symptoms {
				s.Muck()
			}

			// Execute Middleware lifecycle hook: post-muck
			//for _, s := range middlewares {
			//	go func() {
			//    s.Whatever(context...)
			//  }()
			//}

			proxy := &httputil.ReverseProxy{Director: director}
			proxy.ServeHTTP(w, r)
		})
		err := http.ListenAndServe(fmt.Sprintf("%s:%d", m.config.Host, m.config.Port), mux)
		if err != nil {
			log.Println(fmt.Sprintf("ListenAndServe error: ", err))
		}
	}()

	// Block until a signal is received.
	<-sigChan
	log.Println("Shutting down Muxy...")

	//symptom.Teardown()
	for _, s := range symptoms {
		s.Teardown()
	}
}
