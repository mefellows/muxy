package symptom

import (
	"github.com/mefellows/muxy/config"
	"github.com/mefellows/muxy/muxy"
	"log"
	"time"
)

// 50x, 40x etc.

type HttpErrorSymptom struct {
	Delay int `required:"true" default:"2"`
}

const DEFAULT_DELAY = 2 * time.Second

func init() {
	muxy.SymptomFactories.Register(func() (muxy.Symptom, error) {
		return &HttpErrorSymptom{}, nil
	}, "http_error")

}

func (m HttpErrorSymptom) Configure(c *config.RawConfig) error {
	log.Println("HTTP Error Configure()")
	return nil
}

func (m HttpErrorSymptom) Setup() {
	log.Println("HTTP Error Setup()")
}

func (m HttpErrorSymptom) Teardown() {
	log.Println("HTTP Error Teardown()")
}
func (h *HttpErrorSymptom) Muck() {
	delay := time.Duration(h.Delay) * time.Second
	log.Printf("HTTP Error Muck(), delaying for %v seconds\n", delay.Seconds())

	for {
		select {
		case <-time.After(delay):
			return
		}
	}
}
