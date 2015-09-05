package muxy

import (
	"testing"
)

var SymptomMockFunc = func() (Symptom, error) {
	return MockSymptom{}, nil
}

func TestFactoryAll(t *testing.T) {
	SymptomFactories.Register(SymptomMockFunc, "screwything")

	symptoms := SymptomFactories.All()
	for _, sym := range symptoms {
		f, _ := sym()
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

	sym, err := f()

	if err != nil {
		t.Fatalf("Did not expect err: %v", err)
	}

	if sym == nil {
		t.Fatalf("Expected symptom not to be nil")
	}
}
