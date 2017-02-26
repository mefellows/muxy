package symptom

import (
	"time"

	"github.com/mefellows/muxy/log"
	"github.com/mefellows/muxy/muxy"
	"github.com/mefellows/plugo/plugo"
)

// HttpDelaySymptom adds specified delays to HTTP requests
// nolint
type HttpDelaySymptom struct {
	Delay int `required:"true" default:"2"`
}

func init() {
	plugo.PluginFactories.Register(func() (interface{}, error) {
		return &HttpDelaySymptom{}, nil
	}, "http_delay")
}

// Setup sets up the delay plugin
func (m HttpDelaySymptom) Setup() {
	log.Debug("HTTP Delay Setup()")
}

// Teardown shuts down the plugin
func (m HttpDelaySymptom) Teardown() {
	log.Debug("HTTP Delay Teardown()")
}

// HandleEvent takes a proxy event for the proxy to intercept and modify
func (m HttpDelaySymptom) HandleEvent(e muxy.ProxyEvent, ctx *muxy.Context) {
	switch e {
	case muxy.EventPreDispatch:
		m.Muck(ctx)
	}
}

// Muck injects chaos into the system
func (m *HttpDelaySymptom) Muck(ctx *muxy.Context) {
	delay := time.Duration(m.Delay) * time.Second
	log.Debug("HTTP Delay Muck(), delaying for %v seconds\n", delay.Seconds())

	for {
		select {
		case <-time.After(delay):
			return
		}
	}
}
