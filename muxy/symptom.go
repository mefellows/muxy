package muxy

type Symptom interface {
	//
	Setup()
	Muck()
	Teardown()
}
