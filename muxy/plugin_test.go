package muxy

import (
	"github.com/mefellows/muxy/config"
	"testing"
)

var SymptomMockFunc = func(c config.RawConfig) (Symptom, error) {
	return MockSymptom{}, nil
}
var c = config.RawConfig{}

func TestFactoryAll(t *testing.T) {
	SymptomFactories.Register(SymptomMockFunc, "screwything")

	symptoms := SymptomFactories.All()
	for _, sym := range symptoms {
		f, _ := sym(c)
		if _, ok := f.(Symptom); !ok {
			t.Fatalf("must be a Symptom")
		}

	}
}
func TestLookupFactory(t *testing.T) {
	SymptomFactories.Register(SymptomMockFunc, "screwything")

	f, ok := SymptomFactories.Lookup("screwything")

	if !ok {
		t.Fatalf("Expected lookup to be OK")
	}

	sym, err := f(c)

	if err != nil {
		t.Fatalf("Did not expect err: %v", err)
	}

	if sym == nil {
		t.Fatalf("Expected symptom not to be nil")
	}
}
