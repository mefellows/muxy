package muxy

import (
	"log"
	"testing"
)

var MiddlewareMockFunc = func() (interface{}, error) {
	return MockSymptom{}, nil
}

func TestLookupFactory(t *testing.T) {
	PluginFactories.Register(MiddlewareMockFunc, "screwything")

	f, ok := PluginFactories.Lookup("screwything")

	if !ok {
		t.Fatalf("Expected lookup to be OK")
	}

	sym, err := f()

	if err != nil {
		t.Fatalf("Did not expect err: %v", err)
	}

	if sym == nil {
		t.Fatalf("Expected symptom not to be nil")
	}
}

type MockSymptom struct {
	ConfigureError error
	ConfigureCount int
	MuckCount      int
	TeardownCount  int
	SetupCount     int
}

func (m MockSymptom) Setup() {
	log.Println("Mock Setup()")
	m.SetupCount = m.SetupCount + 1
}

func (m MockSymptom) HandleEvent(e ProxyEvent, ctx *Context) {
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
