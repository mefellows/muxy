package symptom

import (
	"log"
	"time"
	"github.com/mefellows/muxy/config"
	"fmt"
)

// 50x, 40x etc.

type HttpErrorSymptom struct{}

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
	log.Println("HTTP Error Muck()")

	for {
		select {
		case <-time.After(5 * time.Second):
			fmt.Printf("Finishing my delay...")
			return
		}
	}
}
