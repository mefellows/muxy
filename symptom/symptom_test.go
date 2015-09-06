package symptom

import (
	"github.com/mefellows/muxy/config"
	"github.com/mefellows/muxy/muxy"
	"log"
)

type MockSymptom struct {
	ConfigureError error
	ConfigureCount int
	MuckCount      int
	TeardownCount  int
	SetupCount     int
}

func (m MockSymptom) Configure(c *config.RawConfig) error {
	log.Println("Mock Configure()")
	m.ConfigureCount = m.ConfigureCount + 1
	return m.ConfigureError
}

func (m MockSymptom) Setup() {
	log.Println("Mock Setup()")
	m.SetupCount = m.SetupCount + 1
}

func (m MockSymptom) HandleEvent(e muxy.ProxyEvent, ctx *muxy.Context) {
	log.Println("Mock HandleEvent()")
}

func (m MockSymptom) Muck() {
	log.Println("Mock Muck()")
	m.MuckCount = m.MuckCount + 1
}

func (m MockSymptom) Teardown() {
	log.Println("Mock Teardown()")
	m.TeardownCount = m.TeardownCount + 1
}
