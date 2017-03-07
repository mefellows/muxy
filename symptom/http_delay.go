package symptom

import (
	"time"

	"github.com/mefellows/muxy/log"
	"github.com/mefellows/muxy/muxy"
	"github.com/mefellows/plugo/plugo"
)

// HTTPDelaySymptom adds specified delays to HTTP requests
// Update docs: these values should be in ms
type HTTPDelaySymptom struct {
	RequestDelay  int `required:"false" mapstructure:"request_delay"`
	ResponseDelay int `required:"false" mapstructure:"response_delay"`
	Delay         int `required:"false" mapstructure:"delay"`
}

func init() {
	plugo.PluginFactories.Register(func() (interface{}, error) {
		return &HTTPDelaySymptom{}, nil
	}, "http_delay")
	plugo.PluginFactories.Register(func() (interface{}, error) {
		return &HTTPDelaySymptom{}, nil
	}, "delay")
}

// Setup sets up the delay plugin
func (m HTTPDelaySymptom) Setup() {
	log.Debug("HTTP Delay Setup()")
}

// Teardown shuts down the plugin
func (m HTTPDelaySymptom) Teardown() {
	log.Debug("HTTP Delay Teardown()")
}

// HandleEvent takes a proxy event for the proxy to intercept and modify
func (m HTTPDelaySymptom) HandleEvent(e muxy.ProxyEvent, ctx *muxy.Context) {
	switch e {
	case muxy.EventPreDispatch:
		if m.RequestDelay != 0 {
			m.Muck(ctx, m.RequestDelay)
		}
	case muxy.EventPostDispatch:
		if m.ResponseDelay != 0 {
			m.Muck(ctx, m.ResponseDelay)
		} else if m.Delay != 0 { // legacy behaviour
			m.Muck(ctx, m.Delay*1000) // convert to ms
		}
	}
}

// Muck injects chaos into the system
func (m *HTTPDelaySymptom) Muck(ctx *muxy.Context, wait int) {
	delay := time.Duration(wait) * time.Millisecond
	log.Debug("HTTP Delay Muck(), delaying for %v seconds\n", delay.Seconds())

	for {
		select {
		case <-time.After(delay):
			return
		}
	}
}
