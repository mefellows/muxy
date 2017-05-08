package symptom

import (
	"time"

	"github.com/mefellows/muxy/log"
	"github.com/mefellows/muxy/muxy"
	"github.com/mefellows/plugo/plugo"
)

// HTTPDelaySymptom adds specified delays to requests Symptom
// Update docs: these values should be in ms
type HTTPDelaySymptom struct {
	RequestDelay  int            `required:"false" mapstructure:"request_delay"`
	ResponseDelay int            `required:"false" mapstructure:"response_delay"`
	Delay         int            `required:"false" mapstructure:"delay"`
	MatchingRules []MatchingRule `required:"false" mapstructure:"matching_rules"`
}

var oneSecondInMillis = 1000

func init() {
	plugo.PluginFactories.Register(func() (interface{}, error) {
		return &HTTPDelaySymptom{}, nil
	}, "http_delay")
	plugo.PluginFactories.Register(func() (interface{}, error) {
		return &HTTPDelaySymptom{}, nil
	}, "delay")
}

// Setup sets up the delay plugin
func (m *HTTPDelaySymptom) Setup() {
	log.Debug("Delay Symptom - Setup()")

	// Add default (catch all) matching rule
	// Only applicable if none supplied
	if len(m.MatchingRules) == 0 {
		m.MatchingRules = []MatchingRule{
			defaultMatchingRule,
		}
	}
}

// Teardown shuts down the plugin
func (m *HTTPDelaySymptom) Teardown() {
	log.Debug("Delay Symptom - Teardown()")
}

// HandleEvent takes a proxy event for the proxy to intercept and modify
func (m *HTTPDelaySymptom) HandleEvent(e muxy.ProxyEvent, ctx *muxy.Context) {
	if MatchSymptoms(m.MatchingRules, *ctx) {
		log.Trace("HTTP Delay Tamperer Hit")

		switch e {
		case muxy.EventPreDispatch:
			if m.RequestDelay != 0 {
				m.Muck(ctx, m.RequestDelay)
			}
		case muxy.EventPostDispatch:
			if m.ResponseDelay != 0 {
				m.Muck(ctx, m.ResponseDelay)
			} else if m.Delay != 0 { // legacy behaviour
				m.Muck(ctx, m.Delay*oneSecondInMillis) // convert to ms
			}
		}
	} else {
		log.Trace("HTTP Delay Tamperer Miss")
	}
}

// Muck injects chaos into the system
func (m *HTTPDelaySymptom) Muck(ctx *muxy.Context, wait int) {
	delay := time.Duration(wait) * time.Millisecond
	log.Debug("Delay Symptom - Muck(), delaying for %v seconds", delay.Seconds())

	for {
		select {
		case <-time.After(delay):
			return
		}
	}
}
