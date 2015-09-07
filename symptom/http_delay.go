package symptom

import (
	"github.com/mefellows/muxy/log"
	"github.com/mefellows/muxy/muxy"
	"github.com/mefellows/plugo/plugo"
	"time"
)

type HttpDelaySymptom struct {
	Delay int `required:"true" default:"2"`
}

const DEFAULT_DELAY = 2 * time.Second

func init() {
	plugo.PluginFactories.Register(func() (interface{}, error) {
		return &HttpDelaySymptom{}, nil
	}, "http_delay")
}

func (m HttpDelaySymptom) Setup() {
	log.Debug("HTTP Delay Setup()")
}

func (m HttpDelaySymptom) Teardown() {
	log.Debug("HTTP Delay Teardown()")
}

func (m HttpDelaySymptom) HandleEvent(e muxy.ProxyEvent, ctx *muxy.Context) {
	switch e {
	case muxy.EVENT_PRE_DISPATCH:
		m.Muck(ctx)
	}
}

func (h *HttpDelaySymptom) Muck(ctx *muxy.Context) {
	delay := time.Duration(h.Delay) * time.Second
	log.Debug("HTTP Delay Muck(), delaying for %v seconds\n", delay.Seconds())

	for {
		select {
		case <-time.After(delay):
			return
		}
	}
}
