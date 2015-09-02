package muxy

import (
	"log"
)

type MockSymptom struct {
	ConfigureError error
	ConfigureCount int
	MuckCount      int
	TeardownCount  int
	SetupCount     int
}

func (m MockSymptom) Configure(c *RawConfig) error {
	log.Println("Mock Configure()")
	m.ConfigureCount = m.ConfigureCount + 1
	return m.ConfigureError
}

func (m MockSymptom) Setup() {
	log.Println("Mock Setup()")
	m.SetupCount = m.SetupCount + 1
}

func (m MockSymptom) Muck() {
	log.Println("Mock Muck()")
	m.MuckCount = m.MuckCount + 1
}

func (m MockSymptom) Teardown() {
	log.Println("Mock Teardown()")
	m.TeardownCount = m.TeardownCount + 1
}
