package muxy

import (
	"github.com/mefellows/muxy/log"
	"github.com/mefellows/plugo/plugo"
	"os"
	"os/signal"
)

type MuxyConfig struct {
	RawConfig  *plugo.RawConfig
	ConfigFile string // Path to YAML Configuration File
}
type PluginConfig struct {
	Name        string
	Description string
	LogLevel    int `default:"2" required:"true" mapstructure:"loglevel"`
	Proxy       []plugo.PluginConfig
	Middleware  []plugo.PluginConfig
}

type Muxy struct {
	config      *MuxyConfig
	middlewares []Middleware
	proxies     []Proxy
}

func New(config *MuxyConfig) *Muxy {
	return &Muxy{config: config}
}

func NewWithDefaultMuxyConfig() *Muxy {
	c := &MuxyConfig{}
	return &Muxy{config: c}
}

func (m *Muxy) Run() {
	m.LoadPlugins()

	// Setup all plugins...
	for _, m := range m.middlewares {
		m.Setup()
	}

	// Interrupt handler
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill)

	// Start proxy
	for _, proxy := range m.proxies {
		go proxy.Proxy()
	}

	// Block until a signal is received.
	<-sigChan
	log.Info("Shutting down Muxy...")

	for _, m := range m.middlewares {
		m.Teardown()
	}
}

func (m *Muxy) LoadPlugins() {
	// Load Configuration
	var err error
	var confLoader *plugo.ConfigLoader
	c := &PluginConfig{}
	if m.config.ConfigFile != "" {
		confLoader = &plugo.ConfigLoader{}
		err = confLoader.LoadFromFile(m.config.ConfigFile, &c)
		if err != nil {
			log.Fatalf("Unable to read configuration file: %s", err.Error())
		}
	} else {
		log.Fatal("No config file provided")
	}

	log.SetLevel(log.LogLevel(c.LogLevel))

	// Load all plugins
	m.middlewares = make([]Middleware, len(c.Middleware))
	plugins := plugo.LoadPluginsWithConfig(confLoader, c.Middleware)
	for i, p := range plugins {
		log.Info("Loading plugin \t" + log.Colorize(log.YELLOW, c.Middleware[i].Name))
		m.middlewares[i] = p.(Middleware)
	}

	m.proxies = make([]Proxy, len(c.Proxy))
	plugins = plugo.LoadPluginsWithConfig(confLoader, c.Proxy)
	for i, p := range plugins {
		log.Info("Loading proxy \t" + log.Colorize(log.YELLOW, c.Proxy[i].Name))
		m.proxies[i] = p.(Proxy)
		m.proxies[i].Setup(m.middlewares)
	}
}
