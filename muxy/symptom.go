package muxy

import (
//"github.com/mefellows/muxy/config"
)

type Symptom interface {
	//Configure(c *config.RawConfig) error // Meet configurable interface
	Setup()
	Muck()
	Teardown()
}

// Register with plugin factory:
//
//func init() {
//  muxy.SymptomFactories.Register(NewSymptomx, "sypmtomname")
//}
//
// Where NewSymptomx is a func that returns a Symptom when called
