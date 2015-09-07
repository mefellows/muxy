package symptom

import (
	"github.com/mefellows/muxy/log"
	"github.com/mefellows/muxy/muxy"
	"github.com/mefellows/plugo/plugo"
	"time"
)

// 50x, 40x etc.

type HttpErrorSymptom struct {
	Delay int `required:"true" default:"2"`
}

const DEFAULT_DELAY = 2 * time.Second

func init() {
	plugo.PluginFactories.Register(func() (interface{}, error) {
		return &HttpErrorSymptom{}, nil
	}, "http_error")
}

func (m HttpErrorSymptom) Setup() {
	log.Debug("HTTP Error Setup()")
}

func (m HttpErrorSymptom) Teardown() {
	log.Debug("HTTP Error Teardown()")
}

func (m HttpErrorSymptom) HandleEvent(e muxy.ProxyEvent, ctx *muxy.Context) {
	switch e {
	case muxy.EVENT_PRE_DISPATCH:
		m.Muck(ctx)
	}
}

func (h *HttpErrorSymptom) Muck(ctx *muxy.Context) {
	delay := time.Duration(h.Delay) * time.Second
	log.Debug("HTTP Error Muck(), delaying for %v seconds\n", delay.Seconds())

	for {
		select {
		case <-time.After(delay):
			return
		}
	}
}
